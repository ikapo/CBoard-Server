const blueButton = document.getElementById('blue');
const greenButton = document.getElementById('green');
const redButton = document.getElementById('red');
const body = document.body;
const theme = localStorage.getItem('theme');

if (theme){
    body.className = theme;
}

blueButton.onclick = () => {
    body.className = "blue";
    localStorage.setItem('theme', 'blue');
}

greenButton.onclick = () => {
    body.className = "green";
    localStorage.setItem('theme', 'green');
}

redButton.onclick = () => {
    body.className = "red";
    localStorage.setItem('theme', 'red');
}

