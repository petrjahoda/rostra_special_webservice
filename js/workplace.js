workplaceBackButton.addEventListener("click", () => {
    operationRow.classList.remove("disabled");
    workplaceRow.classList.add("disabled")
    const select = Metro.getPlugin("#workplace-select", 'select');
    select.data({
    });
    operationSelect.focus()
})

workplaceOkButton.addEventListener("click", () => {
    console.log("Workplace selected: " + workplaceSelect.value);
})
