import React, { useState } from 'react';

function Move(props) {
	const turn = props.turn || 1;
	const white = props.white || "";
	const black = props.black || "";

	return (
		<div className="Move">
			<span className="turn">{turn}.</span>
			<div className="sans">
				<span className="san">{white}</span>
				<span className="san">{black}</span>
			</div>
		</div>
	);
}

function Game(props) {
	const conn = props.conn || null;
	const [game, setGame] = useState({
		moves: [{turn: 1}],
		error: "",
		myTurn: true,
		fen: "",
	});

	if (conn) {
		conn.onmessage = e => {
			game.myTurn = !game.myTurn;
			setGame({...game});
		};
	}
	
	const moveElements = [];
	for (let move of game.moves) {
		moveElements.push(<Move key={move.turn} turn={move.turn} white={move.white} black={move.black} />);
	}

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
	
	const enterMove = e => {
		if (e.charCode !== 13) return;
		const san = e.target.value;
		
		if (!conn || conn.readyState !== WebSocket.OPEN || !game.myTurn) return;

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
					conn.send(JSON.stringify(data));
					makeMove(data.san, data.fen);
				});
			}
		})
		.catch(err => {
			console.log(err);
		});

		e.target.value = "";
	};

	return (
		<div className="Game">
			<div className="moves">
				{ moveElements }
			</div>
			<div className="error">
				{ game.error }
			</div>
			<input type="text" onKeyPress={enterMove} placeholder="Type your move" />
			<a className="fen" href={`https://lichess.org/analysis/standard/${game.fen}`} target="_blank" rel="noopener noreferrer">{game.fen}</a>
		</div>
	);
}

export default Game;
