(function () {
  const form = document.getElementById("form");
  if (!form) return;

  form.addEventListener("submit", sendFormData);
  /**
   * @param {SubmitEvent} e
   */
  function sendFormData(e) {
    e.preventDefault();
    const formData = new FormData(form);
    fetch("http://localhost:8080/api/auth/register", {
      method: "POST",
      headers: {
        // "Content-Type": "application/json",
      },
      body: JSON.stringify({
        email: formData.get("email"),
        nickname: formData.get("nickname"),
        password: formData.get("password"),
        repeatPassword: formData.get("repeat_password"),
      }),
    })
      .then((r) => {
        console.log(r);
      })
      .catch((r) => {
        console.log(r);
      });
  }
})();
