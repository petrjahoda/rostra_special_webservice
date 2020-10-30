transferOrderButton.addEventListener("click", () => {
    console.log("User pressed transfer button")
    let data = {
        WorkplaceCode: sessionStorage.getItem("workplaceCode"),
        UserId: sessionStorage.getItem("userId"),
        OrderInput: sessionStorage.getItem("orderInput"),
        OperationSelect: sessionStorage.getItem("operationValue"),
        OrderId: sessionStorage.getItem("orderId"),
        Nasobnost: sessionStorage.getItem("nasobnost"),
        UserInput: sessionStorage.getItem("userInput"),
        OkCount: sessionStorage.getItem("okCount"),
        NokCount: sessionStorage.getItem("nokCount"),
        NokType: sessionStorage.getItem("nokType")
    };
    fetch("/transfer_order", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                infoRostra.textContent = "Data uložena"
                setTimeout(() => window.location.replace(''), 1500)
            } else {
                infoRostra.textContent = result.TransferOrderError;
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
})