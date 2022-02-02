import React from 'react';

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

function PGNView(props) {
	const moves = props.moves || [{turn: 1}];
	const error = props.error || "";
	const fen = props.fen || "";
	const onMove = props.onMove;
	
	const moveElements = [];
	for (let move of moves) {
		moveElements.push(<Move key={move.turn} turn={move.turn} white={move.white} black={move.black} />);
	}

	const enterMove = e => {
		if (e.charCode !== 13) return;
		if (onMove) {
			onMove(e.target.value);
		}
		e.target.value = "";
	};

	return (
		<div className="PGNView">
			<div className="moves">
				{ moveElements }
			</div>
			<div className="error">
				{ error }
			</div>
			<input type="text" onKeyPress={enterMove} placeholder="Type your move" />
			<div className="fen">
				{ fen }
			</div>
		</div>
	);
}

export default PGNView;
