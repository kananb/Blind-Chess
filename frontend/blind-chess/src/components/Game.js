import React, { useEffect, useRef, useState } from 'react';

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
	const {conn, code} = props;
	const onLeave = props.onLeave || (() => {});
	const [game, setGame] = useState({
		History: [""],
		Error: "",
		Side: "",
		SideToMove: "",
		FEN: "",
		Loser: "",
		WhiteTime: 50,
		BlackTime: 3140,
	});
	const whiteTime = useRef(undefined);
	const blackTime = useRef(undefined);

	/*
	auto scroll dowwn on moves
	focus input on turn
	*/

	// useEffect(() => {
	// 	const countdown = () => {
	// 		setTimeout(() => {
	// 			game.WhiteTime -= 1;
	// 			const min = "0" + Math.floor(game.WhiteTime / 600);
	// 			const sec = "0" + Math.floor(game.WhiteTime / 10 % 60);
	// 			whiteTime.current.innerText = `${min.substring(min.length-2)}:${sec.substring(sec.length-2)}`;
	// 			if (game.WhiteTime > 0) countdown();
	// 		}, 100);
	// 	};
	// 	countdown();
	// });
	
	const moveElements = [];
	let turn = 1;
	if (game.History.length === 0) game.History = [""];
	for (let i = 0; i < game.History.length; i += 2, turn++) {
		moveElements.push(<Move key={turn} turn={turn} white={game.History[i]} black={(i + 1 < game.History.length) ? game.History[i + 1] : ""} />);
	}
	
	const updateGame = (state) => {
		setGame({
			History: state.History || game.History,
			Error: state.Error || "",
			Side: state.Side || game.Side,
			SideToMove: state.SideToMove || game.SideToMove,
			FEN: state.FEN || game.FEN,
			Loser: state.Loser || game.Loser,
		});
	};
	if (conn) {
		conn.onmessage = e => {
			const split = e.data.split("_");
			let msg = {
				cmd: split[0],
				args: split.slice(1),
			};

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

	let info = undefined;
	if (game.FEN) {
		info = <a className="fen" href={`https://lichess.org/analysis/standard/${game.FEN}`} target="_blank" rel="noopener noreferrer">{game.FEN}</a>;
	} else {
		info = <span className="code">room code: { code }</span>;
	}

	let input = undefined;
	if (game.Loser) {
		input = <input type="text" onKeyPress={enterMove} placeholder="Game over" disabled />;
	} else if (!game.FEN) {
		input = <input type="text" onKeyPress={enterMove} placeholder="Waiting for game to start" disabled />;
	} else if (game.Side === game.SideToMove) {
		input = <input className="prompt" type="text" onKeyPress={enterMove} placeholder="Type your move" />;
	} else {
		input = <input type="text" onKeyPress={enterMove} placeholder="Waiting for opponent" disabled />;
	}

	let notification = undefined;
	if (game.Loser) {
		notification = (
			<div className="notification">
				<h3>Game Over</h3>
				{ (game.Loser !== game.Side) ? "You won!!" : "You lost :(" }
			</div>
		);
	}

	return (
		<div className="Game">
			{ notification }
			<div className="clocks">
				<div className="timer active">
					<span ref={whiteTime} className="time"></span>
				</div>
				<div className="timer">
					<span ref={blackTime} className="time"></span>
				</div>
			</div>
			<div className="moves">
				{ moveElements }
			</div>
			<div className="error">
				{ game.Error }
			</div>
			<div className="controls">
				{ input }
				<button className="leave" onClick={() => {
					if (conn) conn.send("QUIT");
					onLeave();
				}}>Leave</button>
			</div>
			{ info }
		</div>
	);
}

export default Game;
