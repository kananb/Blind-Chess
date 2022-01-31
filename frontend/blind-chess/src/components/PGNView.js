import React, { useState } from 'react';

function Move(props) {
	const turn = props.turn || 1;
	const white = props.white || "";
	const black = props.black || "";

	return (
		<div className="Move">
			<span className="turn">{turn}.</span>
			<span className="san">{white}</span>
			<span className="san">{black}</span>
		</div>
	);
}

function PGNView(props) {
	let [moves, setMoves] = useState([{turn: 1}]);
	let [error, setError] = useState("");
	
	const moveElements = [];
	for (let move of moves) {
		moveElements.push(<Move key={move.turn} turn={move.turn} white={move.white} black={move.black} />)
	}

	const enterMove = e => {
		if (e.charCode !== 13) return;
		const san = e.target.value;
		
		if (moves.length === 0 || moves[moves.length-1].black) {
			moves.push({
				turn: moves.length + 1,
				white: san,
			});
		} else if (moves[moves.length-1].white) {
			moves[moves.length-1].black = san;
		} else {
			moves[moves.length-1].white = san;
		}

		e.target.value = "";
		setMoves([...moves]);
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
		</div>
	);
}

export default PGNView;
