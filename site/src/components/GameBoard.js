import React from "react";

const GameBoard = ({ board, ships }) => {
	return (
		<div className="game-board">
			{board?.map((row, rowIndex) => (
				<div key={rowIndex} className="row">
					{row.map((cell, colIndex) => (
						<div key={colIndex} className={`cell ${cell.type}`}>
							{ships?.find(
								(ship) => ship.position.x === colIndex && ship.position.y === rowIndex
							) && <div className="ship">🚢</div>}
						</div>
					))}
				</div>
			))}
		</div>
	);
};

export default GameBoard;

