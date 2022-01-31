import './components/Navbar.css';
import './components/Game.css';
import './components/PGNView.css';
import './components/Chessboard.css';
import './components/Footer.css';
import Footer from './components/Footer';
import Game from './components/Game';
import Navbar from './components/Navbar';

function App() {
  return (
    <div className="App">
		<Navbar />
    	<Game />
		<Footer />
    </div>
  );
}

export default App;
