import React from 'react';
import logo from '../assets/logo_128.png';

function Navbar(props) {
	return (
		<div className="Navbar">
			<img src={logo} alt="Logo" />
			<h1>Blind Chess</h1>
		</div>
	);
}

export default Navbar;
