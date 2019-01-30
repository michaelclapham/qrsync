type Client = {
    qr: string;
    id: string;
    name: string;
    gotoUrl: string;
};

var client: Client;

function createClientOnServer(callback: (client: Client) => any) {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var client = JSON.parse(xhttp.responseText);
            callback(client);
        }
    };
    xhttp.open("GET", "/api/newclient", true);
    xhttp.send();
}

function getClientFromServer(callback: (client: Client) => any) {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var client = JSON.parse(xhttp.responseText);
            callback(client);
        }
    };
    xhttp.open("GET", "/api/client/" + client.id, true);
    xhttp.send();
}

function getQRCode() {
    createClientOnServer((c) => {
        client = c;
        document.getElementById("qr").setAttribute("src", "data:image/png;base64," + client.qr);
        setupNameInput();
    })
}

var debounceTimeout: number;
function setupNameInput() {
    var inputElement = document.getElementById("nameInput") as HTMLInputElement;
    var changeHandler = (event) => {
        client.name = inputElement.value;

        clearTimeout(debounceTimeout);
        if (client.name && client.name.length > 0) {
            debounceTimeout = setTimeout(() => {
                var xhttp = new XMLHttpRequest();
                xhttp.open("POST", `/api/client/${client.id}`, true);
                xhttp.send(JSON.stringify(client));
            }, 1000);
        }
    };
    inputElement.addEventListener("input", changeHandler);
}

getQRCode();

setInterval(() => {
    getClientFromServer((c) => {
        var prevClient = client;
        if (c.gotoUrl && c.gotoUrl.length > 3) {
            location.href = c.gotoUrl;
        }
    });
}, 4000);