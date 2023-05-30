let socket = new WebSocket("ws://"+document.location.host+"/speaker");
let name = "";
let waitForHosts = false;
let isReadyToShoot = false;

const circle = document.getElementById('circle');
const game = document.getElementById("game");

function startGame(){
	document.getElementById("servers").innerHTML = "";
	if (document.getElementById("inpT").value == ""){
		document.getElementById("label").innerText = "Надо ввести ID!!";
	}else{
		name = document.getElementById("inpT").value;
		socket.send(JSON.stringify({type: 1, mess: name}));
		waitForHosts = true;
		document.getElementById("main").style.display = "none";
		document.getElementById("wait").style.display = "flex";
	}
}

function choose(){
	socket.send(JSON.stringify({type: 3}));
	document.getElementById("chooseServer").style.display = "none";
	document.getElementById("wait").style.display = "flex";
}

function chooseSomebody(event){
	socket.send(JSON.stringify({type: 2, mess: event.target.dataset.id}));
	document.getElementById("chooseServer").style.display = "none";
	document.getElementById("wait").style.display = "flex";
}

function killEnemy(event){
	if (isReadyToShoot){
		isReadyToShoot = false;
		document.getElementById("circle").style.background = "red";
		socket.send(JSON.stringify({type: 5, mess: event.target.dataset.id}));
		document.getElementById("game").style.display = "none";
		document.getElementById("win").style.display = "flex";
		document.getElementById("header").style.display = "flex";
		document.getElementById("bullets").style.display = "flex";
		document.getElementById("servers").style.display = "flex";
		setTimeout(function() {
			document.getElementById("win").style.display = "none";
			document.getElementById("wait").style.display = "flex";
			waitForHosts = true;
			document.getElementById("servers").innerHTML = "";
			startGame();
		}, (5000));
	}
}

function readyToshoot(){
	socket.send(JSON.stringify({type: 6}));
	isReadyToShoot = true;
}

socket.onopen = () => {
	console.log("Подключение успешно");
}

socket.onclose = (event) => {
	console.log("Отключение: ", event);
	setTimeout(function() {
		open("./", "_self");
    }, 2500);
}


socket.onerror = (error) => {
	setTimeout(function() {
		open("./", "_self");
    }, 2500);
	console.log("Ошибка: ", error);
}

function addServer(user){
	let div = document.createElement("button");
	div.className = "server";
	div.innerText = user.name;
	div.dataset.id = user.nhash;
	document.getElementById("servers").appendChild(div);
	div.onclick = chooseSomebody;
}

socket.onmessage = (msg) => {
	const ans = JSON.parse(msg.data);
	console.log(ans);
	if (ans.type == 1 && ans.mess != "neok" && waitForHosts){
		document.getElementById("wait").style.display = "none";
		document.getElementById("chooseServer").style.display = "flex";
		if (ans.mess != null) {
			for (let i = 0; i < ans.mess.length; i++){
				addServer(ans.mess[i]);
			}
		}
	} else if (ans.type == 2){
		waitForHosts = false;
		console.log("Противник: "+ans.mess.name, ans.mess.nhash);
		document.getElementById("wait").style.display = "none";
		document.getElementById("header").style.display = "none";
		document.getElementById("bullets").style.display = "none";
		document.getElementById("servers").innerHTML = "";
		document.getElementById("servers").style.display = "none";
		game.style.display = "flex";
		const cleft = ans.xr * 100;
		const ctop = 30 + ans.yr * 35;
		circle.style.left = cleft+"%";
		circle.style.top = ctop+"%";
		resizeCircle();
		circle.innerText = ans.mess.name;
		circle.dataset.id = ans.mess.nhash;
	} else if (ans.type == 3 && ans.mess == "ok"){
		waitForHosts = false;
		console.log("Вы стали хостом");
	} else if (ans.type == 4) {
		let allHosts = document.getElementById("servers").childNodes;
		for (let i = 0; i < document.getElementById("servers").childElementCount; i++){
			if (allHosts[i].dataset.id == ans.mess){
				allHosts[i].remove()
			}
		}
	} else if (ans.type == 5){
		isReadyToShoot = false;
		document.getElementById("circle").style.background = "red";
		document.getElementById("game").style.display = "none";
		document.getElementById("lose").style.display = "flex";
		document.getElementById("header").style.display = "flex";
		document.getElementById("bullets").style.display = "flex";
		document.getElementById("servers").style.display = "flex";
		setTimeout(function() {
			document.getElementById("lose").style.display = "none";
			document.getElementById("wait").style.display = "flex";
			waitForHosts = true;
			startGame();
		}, (5000));
	} else if (ans.type == 6) {
		document.getElementById("circle").style.background = "green";
	}
}

function resizeCircle() {
	const height = game.getBoundingClientRect().height * 0.07;
	circle.style.width = height + 'px';
	circle.style.height = height + 'px';
}

window.addEventListener('resize', resizeCircle);