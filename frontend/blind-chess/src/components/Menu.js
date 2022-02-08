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

			if (msg.cmd === "DENY") {
				setError(msg.args[0]);
			} else if (msg.cmd === "CODE") {
				conn.send("OK");
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
		e.preventDefault();

		let config;
		if (isNaN(duration.current.value)) {
			duration.current.classList.add("error");
			config = null;
		}
		if (isNaN(increment.current.value)) {
			increment.current.classList.add("error");
			config = null;
		}
		
		if (config === null) return;
		config = {
			"Duration": parseInt(duration.current.value || "0")*600,
			"Increment": parseInt(increment.current.value || "0")*10,
			"PlayAs": playAs.current.value,
		};
		try {
			conn.send(`CREATE_${JSON.stringify(config)}`);
		} catch (err) {
			console.error(err);
		}
	};
	const handleTimeChange = e => {
		if (isNaN(e.target.value)) {
			e.target.classList.add("error");
		} else {
			e.target.classList.remove("error");
		}
	};
	
	const code = useRef(null);
	const duration = useRef(null), increment = useRef(null);
	const playAs = useRef(null);
	return (
		<div className="Menu">
			<div className="matchmaking">
				<h2>Start playing</h2>
				<form className="joinForm" onSubmit={handleJoin}>
					<span className="error">{ error }</span>
					<input ref={ code } className="" type="text" placeholder="Enter a room code" />
					<button className="join" type="submit">Join a game</button>
				</form>
				<hr />
				<form className="createForm" onSubmit={handleCreate}>
					<button className="create">Create a game</button>
					<label>Time control</label>
					<div className="timeControl">
						<input ref={duration} onChange={handleTimeChange} type="text" placeholder="min" />
						<input ref={increment} onChange={handleTimeChange} type="text" placeholder="increment" />
					</div>
					<label>Play as</label>
					<select ref={playAs} name="selectSide" id="selectSide">
						<option value="">Random</option>
						<option value="w">White</option>
						<option value="b">Black</option>
					</select>
				</form>
			</div>

			<div className="info">
				<a href="https://www.chess.com/terms/chess-notation" target="_blank" rel="noopener noreferrer">Learn more about chess notation</a>
				<a href="http://" target="_blank" rel="noopener noreferrer">Found a problem?</a>
			</div>
		</div>
	);
}

export default Menu;
