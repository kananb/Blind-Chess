import React, { useEffect, useRef, useState } from 'react';

import useSound from 'use-sound';
import pieceSound from '../assets/piece_sound.mp3'; 

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
	const [playSound] = useSound(pieceSound);
	const {conn, code} = props;
	const onLeave = props.onLeave || (() => {});
	const [game, setGame] = useState({
		History: [""],
		Error: "",
		Side: "",
		SideToMove: "",
		FEN: "",
		Loser: "",
		WhiteClock: 0,
		BlackClock: 0,
	});
	const interval = useRef(0);
	const whiteTime = useRef(undefined);
	const blackTime = useRef(undefined);
	const inputRef = useRef(undefined);
	const moveRef = useRef(undefined);

	const updateClocks = () => {
		let min, sec;
		if (whiteTime.current) {
			min = "0" + Math.floor(Math.abs(game.WhiteClock) / 600);
			sec = "0" + Math.floor(Math.abs(game.WhiteClock) / 10 % 60);
			whiteTime.current.innerText = `${(game.WhiteClock < 0) ? "-" : ""}${min.substring(min.length-2)}:${sec.substring(sec.length-2)}`;
		}
		if (blackTime.current) {
			min = "0" + Math.floor(Math.abs(game.BlackClock) / 600);
			sec = "0" + Math.floor(Math.abs(game.BlackClock) / 10 % 60);
			blackTime.current.innerText = `${(game.BlackClock < 0) ? "-" : ""}${min.substring(min.length-2)}:${sec.substring(sec.length-2)}`;
		}
	};

	useEffect(() => {
		if (!inputRef.current.disabled) inputRef.current.focus();
		moveRef.current.scrollTop = moveRef.current.scrollHeight;
		updateClocks();

		clearInterval(interval.current);
		if (game.Loser || !game.FEN) return;

		interval.current = setInterval((color) => {
			if (color === "w") game.WhiteClock -= 10;
			else game.BlackClock -= 10;

			updateClocks();
		}, 1000, game.SideToMove);
	}, [game]);
	
	const moveElements = [];
	let turn = 1;
	if (game.History.length === 0) game.History = [""];
	for (let i = 0; i < game.History.length; i += 2, turn++) {
		moveElements.push(<Move key={turn} turn={turn} white={game.History[i]} black={(i + 1 < game.History.length) ? game.History[i + 1] : ""} />);
	}
	
	const updateGame = (state) => {
		game.Error = "";
		if (game.SideToMove !== state.SideToMove && game.SideToMove && state.SideToMove) {
			playSound();
		}
		setGame({...Object.assign(game, state)});
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
	const copyPGN = () => {
		const parts = [];
		let turn = 1;
		for (let i = 0; i < game.History.length; i += 2, turn++) {
			parts.push(`${turn}. ${game.History[i]}${(i + 1 < game.History.length) ? " " + game.History[i+1] : ""}`);
		}

		navigator.clipboard.writeText(parts.join(" "));
	};

	let info = undefined;
	if (game.FEN) {
		info = (
			<div className="positionInfo">
				<a className="fen" href={`https://lichess.org/analysis/standard/${game.FEN}`} target="_blank" rel="noopener noreferrer">{game.FEN}</a>
				 &nbsp;&nbsp;--&nbsp;&nbsp;
				<a className="pgn" onClick={copyPGN} alt="Copy PGN data">PGN
				<span className="tooltip">Copy PGN to clipboard</span>
				</a>
			</div>);
	} else {
		info = <span className="code">room code: { code }</span>;
	}

	let placeholder, prompt = false;
	if (game.Loser) placeholder = "Game over";
	else if (!game.FEN) placeholder = "Waiting for game to start";
	else if (game.Side === game.SideToMove) {
		placeholder = "Type your move";
		prompt = true;
	} else placeholder = "Waiting for opponent";

	let notification = undefined;
	if (game.Loser) {
		notification = (
			<div className="notification">
				<h3>Game Over</h3>
				{ (game.Loser === game.Side) ? "You lost :(" : (game.Loser === "-") ? "It's a draw" : "You won!" }
			</div>
		);
	}

	return (
		<div className="Game">
			{ notification }
			<div className="clocks">
				<div className={"timer " + ((game.SideToMove === "w") ? "active" : "")}>
					<span ref={whiteTime} className="time"></span>
				</div>
				<div className={"timer " + ((game.SideToMove === "b") ? "active" : "")}>
					<span ref={blackTime} className="time"></span>
				</div>
			</div>
			<div ref={moveRef} className="moves">
				{ moveElements }
			</div>
			<div className="error">
				{ game.Error }
			</div>
			<div className="controls">
				<input ref={inputRef} className={prompt ? "prompt" : ""} type="text" onKeyPress={enterMove} placeholder={placeholder} disabled={!prompt} />
				<button className="leave" onClick={() => {
					if (!game.FEN || game.Loser || !conn || conn.readyState !== WebSocket.OPEN) {
						if (conn) conn.send("QUIT");
						onLeave();
					} else {
						if (conn) conn.send("RESIGN");
					}
				}}>{(!game.FEN || game.Loser || !conn || conn.readyState !== WebSocket.OPEN) ? "Leave" : "Resign"}</button>
			</div>
			{ info }
		</div>
	);
}

export default Game;
