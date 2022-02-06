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
	const code = props.code || undefined;
	const [game, setGame] = useState({
		History: [""],
		Error: "",
		Side: "",
		SideToMove: "",
		FEN: "",
	});
	
	const moveElements = [];
	let turn = 1;
	for (let i = 0; i < game.History.length; i += 2, turn++) {
		moveElements.push(<Move key={turn} turn={turn} white={game.History[i]} black={(i + 1 < game.History.length) ? game.History[i + 1] : ""} />);
	}
	
	const updateGame = (state) => {
		setGame({
			History: state.History || game.History || [""],
			Error: state.Error || "",
			Side: state.Side || game.Side,
			SideToMove: state.SideToMove || game.SideToMove,
			FEN: state.FEN || game.FEN,
		});
	};
	if (conn) {
		conn.onmessage = e => {
			const split = e.data.split("_");
			let msg = {
				cmd: split[0],
				args: split.slice(1),
			};

			console.log(msg);
			if (msg.cmd === "END") {

			} else if (msg.cmd === "ERROR") {
				updateGame({Error: msg.args[0]});
			} else if (msg.cmd === "STATE") {
				updateGame(JSON.parse(msg.args[0]));
			}
		};
	}
	const enterMove = e => {
		if (e.charCode !== 13) return;
		const san = e.target.value;
		
		if (!conn || conn.readyState !== WebSocket.OPEN) return;
		conn.send(`MOVE_${san}`);

		e.target.value = "";
	};

	const info = (game.FEN) ?
		(
			<a className="fen" href={`https://lichess.org/analysis/standard/${game.FEN}`} target="_blank" rel="noopener noreferrer">{game.FEN}</a>
		) :
		(
			<span className="code">room code: { code }</span>
		);
	return (
		<div className="Game">
			<div className="moves">
				{ moveElements }
			</div>
			<div className="error">
				{ game.Error }
			</div>
			<input type="text" onKeyPress={enterMove} placeholder="Type your move" />
			{ info }
		</div>
	);
}

export default Game;
