import React, { useRef, useEffect } from "react";
import Peer from "simple-peer";

const WebRTCChat = ({ playerId }) => {
	const peers = useRef({});
	const videoRef = useRef(null);

	const handleIncomingCall = (data) => {
		const peer = new Peer({ initiator: false });
		peer.signal(data.signal);
		peer.on("stream", (stream) => {
			videoRef.current.srcObject = stream;
		});
		peers.current[data.callerId] = peer;
	};

	const startCall = (targetPlayerId) => {
		const peer = new Peer({ initiator: true });
		peer.on("signal", (signal) => {
			// Send signal to the target player via signaling server
		});
		peer.on("stream", (stream) => {
			videoRef.current.srcObject = stream;
		});
		peers.current[targetPlayerId] = peer;
	};

	return (
		<div className="webrtc-chat">
			<h3>Audio/Video Chat</h3>
			<video ref={videoRef} autoPlay />
			<button onClick={() => startCall("some-other-player-id")}>Call</button>
		</div>
	);
};

export default WebRTCChat;

