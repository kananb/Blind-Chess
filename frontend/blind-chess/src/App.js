import './components/Navbar.css';
import './components/Home.css';
import './components/Menu.css';
import './components/Game.css';
import './components/Chessboard.css';
import './components/Footer.css';
import Footer from './components/Footer';
import Home from './components/Home';
import Navbar from './components/Navbar';

function App() {
  return (
    <div className="App">
		<Navbar />
    	<Home />
		<Footer />
    </div>
  );
}

export default App;
