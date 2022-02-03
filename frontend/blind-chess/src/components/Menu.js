import React, { useRef, useState } from 'react';

function Menu(props) {
	const conn = props.conn || null;
	const onJoin = props.onJoin || ((id) => {});
	const [error, setError] = useState("");

	if (conn) {
		conn.onmessage = e => {	
			const split = e.data.split("_");
			let msg = {
				cmd: split[0],
				args: split.slice(1),
			};

			console.log(msg);
			if (msg.cmd === "DENY") {
				setError(msg.args[0]);
			} else if (msg.cmd === "CODE") {
				onJoin(msg.args[0]);
			} else {
				console.error(msg);
			}
		};
	}

	const handleJoin = e => {
		try {
			conn.send(`JOIN_${code.current.value}`);
		} catch (err) {
			console.error(err);
		}

		code.current.select();
		e.preventDefault();
	};
	const handleCreate = e => {
		try {
			conn.send("CREATE");
		} catch (err) {
			console.error(err);
		}
	};
	
	const code = useRef(null);
	return (
		<div className="Menu">
			<div className="matchmaking">
				<h2>Start playing</h2>
				<form className="joinForm" onSubmit={handleJoin}>
					<span className="error">{ error }</span>
					<input ref={ code } className="" type="text" placeholder="Enter a room code" />
					<button className="join" type="submit">Join a game</button>
				</form>
				<button className="create" onClick={handleCreate}>Create a game</button>
			</div>

			<div className="info">
				<a href="https://www.chess.com/terms/chess-notation" target="_blank" rel="noopener noreferrer">Learn more about chess notation</a>
				<a href="http://" target="_blank" rel="noopener noreferrer">Found a problem?</a>
			</div>
		</div>
	);
}

export default Menu;
