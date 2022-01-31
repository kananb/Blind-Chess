import React from 'react';
import PGNView from './PGNView';
import Chessboard from './Chessboard';

function Home(props) {
	return (
		<div className="Home">
			<PGNView />
			<Chessboard orientation="white" />
		</div>
	);
}

export default Home;
