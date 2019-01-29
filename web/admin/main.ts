// import { BrowserQRCodeReader, VideoInputDevice } from 'zxing-typescript/src/browser/BrowserQRCodeReader'
declare var ZXing;

type Client = {
    qr: string;
    id: string;
    name: string;
};

var clientList: Client[] = [];
var videoElement = document.getElementById("videoElement") as HTMLVideoElement;
var outputText = document.getElementById("outputText") as HTMLParagraphElement;
var clientListElement = document.getElementById("clientList") as HTMLUListElement;
const codeReader = new ZXing.BrowserQRCodeReader();

function tryScan() {
    codeReader.getVideoInputDevices().then((devices) => {
        try {
            const firstDevice = devices[0];
            var found = false;
            codeReader.decodeFromInputVideoDevice(firstDevice.deviceId, videoElement).then((result) => {
                console.log("Result?");
                var qrText = result.getText();
                const newClient: Client = {
                    id: qrText,
                    name: "",
                    qr: ""
                };
                clientList.push(newClient);
                updateClientListElement();
                setTimeout(() => {
                    outputText.textContent = "";
                    tryScan();
                }, 1000);
            });
        } catch (ex) {
            console.error("Something went wrong!", ex);
        }
    });
}

function updateClientListElement() {
    clientListElement.innerHTML = clientList.map((client) => `<li>${client.id}</li>`).join("\n");
}

tryScan();