import snarkdown from 'snarkdown';
let $ = document.querySelector.bind(document);

// Not using #in if running js
$('#in').name = "empty"
$('#js-text').name = "text"

function run() {
    let html = snarkdown($('#in').value);
    $('#out').innerHTML = html;
    $('#js-text').value = html;
}

$('#in').oninput = run;

run();
