endOrderButton.addEventListener("click", () => {
    console.log("User pressed end button")
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
        NokType: sessionStorage.getItem("nokType"),
        RadioSelect: sessionStorage.getItem("radio"),
        TypZdrojeZapsi: sessionStorage.getItem("typZdrojeZapsi")
    };
    fetch("/end_order", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                infoRostra.textContent = "Data uložena"
                setTimeout(() => window.location.replace(''), 3000)
            } else {
                infoRostra.textContent = result.EndOrderError;
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
})