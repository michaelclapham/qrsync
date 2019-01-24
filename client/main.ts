var xhttp = new XMLHttpRequest();
xhttp.onreadystatechange = function () {
    if (this.readyState == 4 && this.status == 200) {
        var bodyJson = JSON.parse(xhttp.responseText);
        document.getElementById("qr").setAttribute("src", "data:image/png;base64," + bodyJson.qr);
    }
};
xhttp.open("GET", "/newclient", true);
xhttp.send();