// import { BrowserQRCodeReader, VideoInputDevice } from 'zxing-typescript/src/browser/BrowserQRCodeReader'
declare var ZXing;

type Client = {
    qr: string;
    id: string;
    name: string;
    gotoUrl: string;
};

var adminId: string = "a" + Math.round(Math.random() * 100000);

var clientList: Client[] = [];
var videoElement = document.getElementById("videoElement") as HTMLVideoElement;
var outputText = document.getElementById("outputText") as HTMLParagraphElement;
var clientListElement = document.getElementById("clientList") as HTMLUListElement;
var urlInputElement = document.getElementById("urlInput") as HTMLInputElement;

const codeReader = new ZXing.BrowserQRCodeReader();

var debounceTimeout: number;
function setupUrlInput() {
    var changeHandler = (event) => {
        clientList.forEach((client) => {
            client.gotoUrl = urlInputElement.value;
        });

        clearTimeout(debounceTimeout);
        if (urlInputElement.value && urlInputElement.value.length > 0) {
            debounceTimeout = setTimeout(() => {
                postAllClients();
            }, 1000);
        }
    };
    urlInputElement.addEventListener("input", changeHandler);
}

setupUrlInput();

function tryScan() {
    codeReader.getVideoInputDevices().then((devices) => {
        try {
            const firstDevice = devices[0];
            var found = false;
            codeReader.decodeFromInputVideoDevice(firstDevice.deviceId, videoElement).then((result) => {
                console.log("Result?");
                var qrText = result.getText();
                getClientFromServer(qrText, (client) => {
                    clientList.push(client);
                    updateClientListElement();
                    setTimeout(() => {
                        outputText.textContent = "";
                        tryScan();
                    }, 1000);
                });
            });
        } catch (ex) {
            console.error("Something went wrong!", ex);
        }
    });
}

function updateClientListElement() {
    clientListElement.innerHTML = clientList.map((client) => `<li>${client.id}</li>`).join("\n");
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
    xhttp.send(JSON.stringify(clientList));
}

tryScan();