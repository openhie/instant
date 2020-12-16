function Setup() {
    var url = "http://localhost:27517/index.html";
    window.location = url;
    fetch("http://localhost:27517/disclaimer")
    fetch("http://localhost:27517/setup")
}


function accept() {
    // document.getElementById("accept").innerHTML = "Hello World";
    fetch('http://localhost:27517/api/makedisclaimer')
        .then(
            function (response) {
                if (response.status !== 200) {
                    console.log('Looks like there was a problem. Status Code: ' +
                        response.status);
                    return;
                }

                // Examine the text in the response
                response.json().then(function (data) {
                    console.log(data);
                });
            }
        )
        .catch(function (err) {
            console.log('Fetch Error :-S', err);
        });
}

var source = new EventSource('http://localhost:27517/events?stream=messages');
source.onmessage = function (e) {
    document.getElementById('console').innerHTML += e.data + '<br>';
};

function Setup() {
    fetch("http://localhost:27517/setup")
}

function DebugDocker() {
    fetch("http://localhost:27517/debugdocker")
}

function DebugKubernetes() {
    fetch("http://localhost:27517/debugkubernetes")
}

function ListDocker() {
    fetch("http://localhost:27517/listdocker")
}

function ComposeUpCoreDOD() {
    fetch("http://localhost:27517/composeupcoredod")
}
