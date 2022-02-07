import React, { useState } from 'react';

function Chessboard(props) {
	const [orientation, setOrientation] = useState("white");
	const board = [];

	const ranks = ["a", "b", "c", "d", "e", "f", "g", "h"];
	for (let r = 0; r < 8; r++) {
		const squares = [];
		for (let c = 0; c < 8; c++) {
			const off = r&1 ? 1 : 0;
			const color = (r*8+c+off)&1 ? "dark" : "light";

			const file = r === 7 ? (orientation === "white" ? ranks[c] : ranks[7-c]) : "";
			const rank = c === 0 ? (orientation === "white" ? (8-r) : r+1) : "";
			squares.push(
				<td key={c} className={ "square "+color }>
					<div className="rank">
						{ rank }
					</div>
					<div className="file">
						{ file }
					</div>
				</td>
			);
		}

		board.push(
			<tr key={r}>
				{ squares }
			</tr>
		);
	}

	const flipBoard = () => {
		if (orientation === "white") {
			setOrientation("black");
		} else {
			setOrientation("white");
		}
	};

	return (
		<div className="Chessboard">
			<table>
				<tbody>
					{ board }
				</tbody>
			</table>
			<div className="controls">
				<button className="flip" onClick={flipBoard}>flip</button>
			</div>
		</div>
	);
}

export default Chessboard;
