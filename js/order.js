orderBackButton.addEventListener("click", () => {
    orderInput.placeholder = ""
    orderInput.value = ""
    userInput.value = ""
    userInput.placeholder = "Zadejte osobní číslo"
    userInput.focus()
    userRow.classList.remove("disabled");
    orderRow.classList.add("disabled")
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
    console.log("Order entered: " + orderInput.value);
    let data = {OrderInput: orderInput.value};
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
                operationSelect.placeholder = "Zadejte číslo operace výrobního příkazu";
                sessionStorage.setItem("orderId", result.OrderId)
                sessionStorage.setItem("orderName", result.OrderName)
                sessionStorage.setItem("orderInput", result.OrderInput)
                let operace = {};
                for (operation of result.Operations) {
                    operace[operation.Operce] = operation.Operce + ": " + operation.Pracoviste + " [" + operation.PracovistePopis + "]"
                }
                const select = Metro.getPlugin("#operation-select", 'select');
                select.data({
                    "Načtené operace": operace
                });
                operationSelect.focus();
            } else {
                console.log(result.OrderError)
                orderInput.placeholder = result.OrderError
                orderInput.value = ""
            }
        });
    }).catch((error) => {
        errorInfoPanel.textContent = error.toString()
    });
}