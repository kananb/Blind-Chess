import React, { useEffect, useState } from 'react';
import Menu from './Menu';
import Game from './Game';
import Chessboard from './Chessboard';

function Home(props) {
	const [conn, setConn] = useState(null);
	const [code, setCode] = useState(undefined);

	useEffect(() => {
		let socket = new WebSocket(`ws://${window.location.host}/game`);

		socket.onopen = () => {
			console.log("Opened websocket connection");
			setConn(socket);
		};

		socket.onclose = () => {

		};

		socket.onerror = () => {

		};

		return () => {
			socket.close();
		};
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
