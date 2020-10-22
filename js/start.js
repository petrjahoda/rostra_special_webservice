startOrderButton.addEventListener("click", () => {
    console.log("User pressed start button")
    let data = {
        WorkplaceCode: sessionStorage.getItem("workplaceCode"),
        UserId: sessionStorage.getItem("userId"),
        OrderInput: sessionStorage.getItem("orderInput"),
        OperationSelect: sessionStorage.getItem("operationValue"),
        RadioSelect: sessionStorage.getItem("radio"),
        ProductId: sessionStorage.getItem("productId"),
        OrderId: sessionStorage.getItem("orderId"),
        TypZdrojeZapsi: sessionStorage.getItem("typZdrojeZapsi"),
        UserInput: sessionStorage.getItem("userInput"),
        Nasobnost: sessionStorage.getItem("nasobnost")
    };
    fetch("/start_order", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                infoRostra.textContent = "Data uložena"
                setTimeout(() => window.location.replace(''), 3000)
            } else {
                infoError.textContent = result.StartOrderError;
            }
        });
    }).catch((error) => {
        infoError.textContent = error.toString()
    });
})