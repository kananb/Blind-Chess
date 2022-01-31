import React, { useEffect, useState } from 'react';
import PGNView from './PGNView';
import Chessboard from './Chessboard';

function Game(props) {
	let socket;
	let [game, setGame] = useState({
		moves: [{turn: 1}],
		error: "",
		myTurn: true,
		fen: "",
	});

	const makeMove = (san, fen, error="") => {
		game.error = error;

		if (san && fen) {
			const last = game.moves.length - 1;
			if (game.moves[last].black) {
				game.moves.push({turn: game.moves.length+1, white: san});
			} else if (game.moves[last].white) {
				game.moves[last].black = san;
			} else {
				game.moves[last].white = san;
			}

			game.fen = fen;
			game.myTurn = !game.myTurn;
		}

		setGame({...game});
	};

	useEffect(() => {
		try {
			socket = new WebSocket("ws://localhost:8080/game");
		} catch (err) {
			console.log("Couldn't establish connection to server");
			return;
		}

		socket.onopen = () => {
			console.log("Opened websocket connection");
		};

		socket.onmessage = e => {
			game.myTurn = true;
			setGame({...game});
		};
	});

	const onMove = san => {
		if (!socket || socket.readyState !== WebSocket.OPEN || !game.myTurn) return;

		fetch(`http://localhost:8080/api/moves/${san}?fen=${game.fen}`, {
			mode: "cors",
			method: "POST",
		})
		.then(res => {
			if (!res.ok) {
				res.text()
				.then(data => {
					makeMove(null, null, data);
				});
			} else {
				res.json()
				.then(data => {
					socket.send(JSON.stringify({
						san: san,
						fen: data.fen,
					}));
					makeMove(san, data.fen);
				});
			}
		})
		.catch(err => {
			console.log(err);
		});
	};

	return (
		<div className="Game">
			<PGNView moves={[...game.moves]} error={game.error} fen={game.fen} onMove={onMove} />
			<Chessboard orientation="white" />
		</div>
	);
}

export default Game;
