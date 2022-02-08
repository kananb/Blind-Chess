import React, { useEffect, useRef, useState } from 'react';
import Menu from './Menu';
import Game from './Game';
import Chessboard from './Chessboard';

function Home(props) {
	const [conn, setConn] = useState(null);
	const [code, setCode] = useState(undefined);
	const timeout = useRef(250);

	useEffect(() => {
		(function connect() {
			let socket = new WebSocket(`ws${(window.location.protocol === "https:") ? "s" : ""}://${window.location.host}/game`);
			let connectInterval;
	
			socket.onopen = () => {
				console.log("Socket connection established.");
				timeout.current = 250;
				clearTimeout(connectInterval);
				setConn(socket);
			};
	
			socket.onclose = () => {
				timeout.current += timeout.current;
				connectInterval = setTimeout(() => {
					if (!conn || conn.readyState === WebSocket.CLOSED) connect();
				}, Math.min(timeout.current, 10000));
			};
	
			socket.onerror = () => {
				console.error("Socket encountered an error, retrying connection.");
				socket.close();
			};
		})();
	}, []);

	const handleJoin = id => {
		setCode(id);
	};
	const handleLeave = () => {
		setCode(undefined);
	};

	let interact = undefined;
	if (code) {
		interact = <Game conn={conn} code={code} onLeave={handleLeave} />;
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
