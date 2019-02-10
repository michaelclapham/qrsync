type Client = {
    qr: string;
    id: string;
    name: string;
    gotoUrl: string;
};

var client: Client;

function getClientFromServer(callback: (client: Client) => any) {
    const xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            var responseClient = JSON.parse(xhttp.responseText);
            callback(responseClient);
        }
    };
    let url = "/api/client";
    if (client && client.id) {
        url += "/" + client.id;
    }
    xhttp.open("GET", url, true);
    xhttp.send();
}

function getQRCode() {
    getClientFromServer((c) => {
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
        if (c.gotoUrl && c.gotoUrl.length > 3) {
            location.href = c.gotoUrl;
        }
    });
}, 4000);