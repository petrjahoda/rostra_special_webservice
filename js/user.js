userInput.addEventListener("keyup", function (event) {
    if (event.code === "Enter") {
        processUserInput();
    }
});


userOkButton.addEventListener("click", () => {
    processUserInput();
})

function processUserInput() {
    console.log("User entered: " + userInput.value);
    let data = {UserInput: userInput.value};
    fetch("/check_user_input", {
        method: "POST",
        body: JSON.stringify(data)
    }).then((response) => {
        response.text().then(function (data) {
            let result = JSON.parse(data);
            if (result.Result === "ok") {
                userInput.value = result.UserName;
                userRow.classList.add("disabled");
                orderRow.classList.remove("disabled");
                orderInput.placeholder = "Zadejte číslo výrobního příkazu";
                sessionStorage.setItem("userId", result.UserId)
                sessionStorage.setItem("userName", result.UserName)
                sessionStorage.setItem("userInput", result.UserInput)
                orderInput.focus();
            } else {
                console.log(result.UserError)
                userInput.placeholder = result.UserError
                userInput.value = ""
            }
        });
    }).catch((error) => {
        errorInfoPanel.textContent = error.toString()
    });
}