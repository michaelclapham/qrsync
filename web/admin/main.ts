// import { BrowserQRCodeReader, VideoInputDevice } from 'zxing-typescript/src/browser/BrowserQRCodeReader'
declare var ZXing;

type Client = {
    qr: string;
    id: string;
    name: string;
    gotoUrl: string;
};

var adminId: string = "a" + Math.round(Math.random() * 100000);

var clientMap: Record<string, Client> = {};
var videoElement = document.getElementById("videoElement") as HTMLVideoElement;
var outputText = document.getElementById("outputText") as HTMLParagraphElement;
var clientListElement = document.getElementById("clientList") as HTMLUListElement;
var urlInputElement = document.getElementById("urlInput") as HTMLInputElement;
var submitButton = document.getElementById("submitButton") as HTMLButtonElement;

const codeReader = new ZXing.BrowserQRCodeReader();

var debounceTimeout: number;
submitButton.onclick = () => {
    clearTimeout(debounceTimeout);
    if (urlInputElement.value && urlInputElement.value.length > 0) {
        debounceTimeout = setTimeout(() => {
            outputText.innerHTML = urlInputElement.value;
            Object.keys(clientMap).forEach((id) => {
                clientMap[id].gotoUrl = urlInputElement.value;
            });
            postAllClients();
        }, 1000);
    }
};

function tryScan() {
    codeReader.getVideoInputDevices().then((devices) => {
        try {
            outputText.textContent = "Number of video devices " + devices.length;
            const firstDevice = devices[devices.length - 1];
            codeReader.decodeFromInputVideoDevice(firstDevice.deviceId, videoElement).then((result) => {
                console.log("Result?");
                var qrText = result.getText();
                setTimeout(() => {
                    outputText.textContent = "";
                    tryScan();
                }, 1000);
                if (!clientMap[qrText]) {
                    getClientFromServer(qrText, (client) => {
                        clientMap[client.id] = client;
                        updateClientListElement();
                    });
                }
            });
        } catch (ex) {
            console.error("Something went wrong!", ex);
        }
    });
}

function updateClientListElement() {
    clientListElement.innerHTML = Object.keys(clientMap).map((clientId) => `<li>${clientId}</li>`).join("\n");
}

function getClientFromServer(id: string, callback: (c: Client) => any) {
    var xhttp = new XMLHttpRequest();
    xhttp.open("GET", `/api/client/${id}`, true);
    xhttp.onreadystatechange = function () {
        if (xhttp.readyState == 4 && xhttp.status == 200) {
            var client = JSON.parse(xhttp.responseText);
            callback(client);
        }
    };
    xhttp.send();
}

function postAllClients() {
    var xhttp = new XMLHttpRequest();
    xhttp.open("POST", `/api/clients`, true);
    xhttp.send(JSON.stringify(Object.keys(clientMap).map((id) => clientMap[id])));
}

tryScan();