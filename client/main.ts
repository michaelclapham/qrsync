type Client = {
    qr: string;
    id: string;
    name: string;
};

var client: Client;

function getQRCode() {
    var xhttp = new XMLHttpRequest();
    xhttp.onreadystatechange = function () {
        if (this.readyState == 4 && this.status == 200) {
            client = JSON.parse(xhttp.responseText);
            document.getElementById("qr").setAttribute("src", "data:image/png;base64," + client.qr);
            setupNameInput();
        }
    };
    xhttp.open("GET", "/newclient", true);
    xhttp.send();
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
                xhttp.open("POST", `/client/${client.id}`, true);
                xhttp.send(JSON.stringify(client));
            }, 1000);
        }
    };
    inputElement.addEventListener("input", changeHandler);
}

getQRCode();
