orderBackButton.addEventListener("click", () => {
    orderInput.placeholder = ""
    orderInput.value = ""
    userInput.value = ""
    userInput.placeholder = "Zadejte osobní číslo"
    userInputCell.style.pointerEvents = "auto"
    orderRow.classList.add("disabled")
    userOkButton.disabled = false;
    orderOkButton.disabled = true;
    orderBackButton.disabled = true;
    infoRostra.textContent = ""
    infoOrderPriznakSeriovaVyroba.textContent = ""
    infoOrderInput.textContent = ""
    infoOrderName.textContent = ""
    infoOrderId.textContent = ""
    const select = Metro.getPlugin("table", 'select');
    select.data({});
    userInput.focus()
})

orderOkButton.addEventListener("click", () => {
    console.log("order button clicked")
    processOrderInput();

})

orderInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processOrderInput();
    }
});

function processOrderInput() {
    console.log("Order value: " + orderInput.value);
    let data = {OrderInput: orderInput.value, UserInput: sessionStorage.getItem("userInput")};
    fetch("/check_order_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                orderInput.value = result.OrderName;
                orderRow.classList.add("disabled");
                operationRow.classList.remove("disabled");
                orderOkButton.disabled = true;
                orderBackButton.disabled = true;
                operationOkButton.disabled = false;
                operationBackButton.disabled = false;
                operationSelect.placeholder = "Zadejte číslo operace výrobního příkazu";
                sessionStorage.setItem("orderId", result.OrderId)
                infoOrderId.textContent = result.OrderId
                sessionStorage.setItem("orderName", result.OrderName)
                infoOrderName.textContent = result.OrderName
                sessionStorage.setItem("orderInput", result.OrderInput)
                infoOrderInput.textContent = result.OrderInput
                sessionStorage.setItem("priznakSeriovaVyroba", result.PriznakSeriovaVyroba)
                sessionStorage.setItem("productId", result.ProductId)
                infoOrderPriznakSeriovaVyroba.textContent = result.PriznakSeriovaVyroba
                let tableData = {};
                for (let operation of result.Operations) {
                    tableData[operation.Operace] = operation.Operace + ": " + operation.Pracoviste + " [" + operation.PracovistePopis + "]"
                }
                const select = Metro.getPlugin("#operation-select", 'select');
                select.data({
                    "Načtené operace": tableData
                });
                infoRostra.textContent = ""
                operationSelect.focus();
            } else {
                infoRostra.textContent = result.OrderError
                orderInput.placeholder = result.OrderError
                orderInput.value = ""
            }
        });
    }).catch((error) => {
        infoRostra.textContent = error.toString()
    });
}