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
                userInputCell.style.pointerEvents = "none"
                orderInput.focus();
                userOkButton.disabled = true
                orderBackButton.disabled = false
                orderOkButton.disabled = false
                orderRow.classList.remove("disabled");
                resetButton.disable = false;
                resetButton.classList.remove("disabled");
                orderInput.placeholder = "Zadejte číslo výrobního příkazu";
                sessionStorage.setItem("userId", result.UserId)
                infoUserId.textContent = result.UserId
                sessionStorage.setItem("userName", result.UserName)
                infoUserName.textContent = result.UserName
                sessionStorage.setItem("userInput", result.UserInput)
                infoUserInput.textContent = result.UserInput
                infoError.textContent = ""
                infoRostra.textContent = ""
                infoError.textContent = ""
                displayTable(result.TableData)
                orderInput.focus();
            } else {
                infoError.textContent = result.UserError;
                userInput.placeholder = result.UserError;
                userInput.value = ""
            }
        });
    }).catch((error) => {
        infoError.textContent = error.toString()
    });
}