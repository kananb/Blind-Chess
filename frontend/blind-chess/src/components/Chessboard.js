import React from 'react';

function Chessboard(props) {
	const board = [];

	const ranks = ["a", "b", "c", "d", "e", "f", "g", "h"];
	for (let r = 0; r < 8; r++) {
		const squares = [];
		for (let c = 0; c < 8; c++) {
			const off = r&1 ? 1 : 0;
			const color = (r*8+c+off)&1 ? "dark" : "light";

			const file = r === 7 ? ranks[c] : "";
			const rank = c === 0 ? (8-r) : "";
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

	return (
		<div className="Chessboard">
			<table>
				<tbody>
					{ board }
				</tbody>
			</table>
		</div>
	);
}

export default Chessboard;
