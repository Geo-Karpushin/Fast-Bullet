const messagesCount = getRandomArbitrary(5, 25);
const timeSpread = 5.5;
const scaleSpread = 2;
const rightSpread = 100;
const topSpread = 70;


window.onload = function(){
	const bodyDOM = document.getElementById("bullets");

	for (let i = 0; i <= messagesCount; i++) {
		let tdelay=getRandomArbitrary(1.5, timeSpread);
		let tscale=getRandomArbitrary(1/scaleSpread, scaleSpread-1.5);
		let ttop=Math.round(getRandomArbitrary(-topSpread,topSpread));
		let tright=Math.round(getRandomArbitrary(-rightSpread+10,rightSpread-10));
		bodyDOM.innerHTML+=`<div class='light' style='-webkit-animation: floatUp `+tdelay+`s infinite linear;-moz-animation: floatUp `+tdelay+`s infinite linear;-o-animation: floatUp `+tdelay+`s infinite linear;animation: floatUp `+tdelay+`s infinite linear;-webkit-transform: scale(`+tscale+`);-moz-transform: scale(`+tscale+`);-o-transform: scale(`+tscale+`);transform: scale(`+tscale+`);top: `+ttop+`%;right: `+tright+`%;'></div>`;
	}
}

function getRandomArbitrary(min, max) {
  return Math.random() * (max - min) + min;
}