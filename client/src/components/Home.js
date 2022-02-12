import React, { useEffect, useRef, useState } from 'react';
import Menu from './Menu';
import Game from './Game';
import Chessboard from './Chessboard';

function Home(props) {
	const [conn, setConn] = useState(null);
	const [joined, setJoined] = useState(false);
	const info = useRef({});
	const timeout = useRef(250);

	useEffect(() => {
		(function connect() {
			let socket = new WebSocket(`ws${(window.location.protocol === "https:") ? "s" : ""}://${window.location.host}/game`);
			let connectTimeout;
	
			socket.onopen = () => {
				console.log("Socket connection established.");
				timeout.current = 250;
				clearTimeout(connectTimeout);

				if (info.current.code) {
					socket.send(`JOIN_${info.current.code}_${info.current.id}`);
				}
				setConn(socket);
			};
	
			socket.onclose = () => {
				console.log("Socket closed, retrying connection.");
				timeout.current += timeout.current;
				connectTimeout = setTimeout(() => {
					if (!conn || conn.readyState === WebSocket.CLOSED) connect();
				}, Math.min(timeout.current, 2000));
			};
	
			socket.onerror = () => {
				console.error("Socket encountered an error, closing connection.");
				socket.close();
			};
		})();
	}, []);

	const handleJoin = (code, id) => {
		info.current = {
			code: code,
			id: id,
		};
		setJoined(true);
	};
	const handleLeave = () => {
		info.current = {};
		setJoined(false);
	};

	let interact = undefined;
	if (joined) {
		interact = <Game conn={conn} code={info.current.code} onLeave={handleLeave} />;
	} else {
		interact = <Menu conn={conn} onJoin={handleJoin} />;
	}

	return (
		<div className="Home">
			<div className="container">
				{ interact }
			</div>
			<Chessboard orientation="white" />
		</div>
	);
}

export default Home;
