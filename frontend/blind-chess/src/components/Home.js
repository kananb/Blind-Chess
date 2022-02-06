import React, { useEffect, useState } from 'react';
import Menu from './Menu';
import Game from './Game';
import Chessboard from './Chessboard';

function Home(props) {
	const [conn, setConn] = useState(null);
	const [code, setCode] = useState(undefined);

	useEffect(() => {
		let socket = new WebSocket("ws://localhost:80/game");

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

	const interact = (code) ?
		(
			<Game conn={conn} code={code} onLeave={handleLeave} />
		) :
		(
			<Menu conn={conn} onJoin={handleJoin} />
		);

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
