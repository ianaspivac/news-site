import "regenerator-runtime/runtime";
const axios = require("axios").default;

async function App() {
  const BASE_URL = "http://localhost:8080";

  let cardList = [];

  let accountFlag = false;

  const getNews = async () => {
    try {
      const response = await axios.get(`${BASE_URL}/v1/posts/`);

      cardList = response.data;
      if (
        !localStorage.getItem("token") ||
        localStorage.getItem("token") === undefined
      ) {
        cardList = cardList.filter((post) => !post.is_protected);
      }
      document
        .getElementById("post-list")
        .insertAdjacentHTML("beforeend", listTemplate());
    } catch (errors) {
      console.error(errors);
    }
  };

  const login = async () => {
    let mail = document.getElementById("email").value;
    let password = document.getElementById("password").value;
    try {
      const response = await axios.post(`${BASE_URL}/v1/user/login`, {
        mail,
        password,
      });

      const user = response.data;
      if (user.role) localStorage.setItem("role", user.role);
      localStorage.setItem("token", user.token);
      window.location.reload(true);
    } catch (errors) {
      document
        .getElementById("form-footer")
        .insertAdjacentHTML("beforebegin", showError(errors));
      console.error(errors);
    }
  };

  const register = async () => {
    let mail = document.getElementById("email").value;
    let password = document.getElementById("password").value;
    try {
      const response = await axios.post(`${BASE_URL}/v1/user/register`, {
        mail,
        password,
      });

      const user = response.data;

      document.getElementById("form-register").remove();
      document.getElementById("account").value = "Logout";

      if (user.role) localStorage.setItem("role", user.role);

      try {
        const response = await axios.post(`${BASE_URL}/v1/user/login`, {
          mail,
          password,
        });

        const user = response.data;

        localStorage.setItem("token", user.token);

        window.location.reload(true);
      } catch (errors) {
        document;
        console.error(errors);
      }
    } catch (error) {
      document
        .getElementById("form-footer")
        .insertAdjacentHTML(
          "beforebegin",
          showError(error.response.data.error)
        );
    }
  };

  const createEl = (name, id = null) => {
    let el = document.createElement(name);
    if (id) {
      el.setAttribute("id", id);
    }
    return el;
  };

  const toBase64 = (file) =>
    new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = (error) => reject(error);
    });

  let postsNr = cardList.length;

  const listTemplate = () => {
    return cardList.map((card, index) => postCard({ card, index })).join("");
  };

  const showContent = (contentId, id) => {
    let container = createEl("div", `container-${contentId}`);
    let p = createEl("p", contentId);
    let contentText = document.createTextNode(cardList[id].content);

    p.appendChild(contentText);
    container.appendChild(p);
    document.getElementById(`summary-container-${id}`).after(container);
  };

  const addEditButton = (id) => {
    let div = createEl("div", `edit-container-${id}`);
    let btn = createEl("button", `edit-${id}`);
    btn.innerHTML = "Edit";
    div.appendChild(btn);
    document.getElementById(`card-${id}`).appendChild(div);
  };

  const toggleFull = (id) => {
    let contentId = `content-${id}`;
    if (!document.getElementById(contentId)) {
      showContent(contentId, id);
      document.getElementById(`toggle-${id}`).value = "See less";
      if (localStorage.getItem("role") === "EDITOR") addEditButton(id);
    } else {
      document.getElementById(`container-${contentId}`).remove();
      document.getElementById(`edit-container-${id}`).remove();
      document.getElementById(`toggle-${id}`).value = "See more";
    }
  };

  const saveForm = async (id) => {
    try {
      const response = await axios.patch(
        `${BASE_URL}/v1/posts/${cardList[id].uuid}`,
        {
          headline: document.getElementById(`new-headline-${id}`).value,
          summary: document.getElementById(`new-summary-${id}`).value,
          content: document.getElementById(`new-content-${id}`).value,
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
            "Content-Type": "application/json",
            Accept: "application/json",
          },
        }
      );
      window.location.reload(true);
    } catch (errors) {
      console.error(errors);
    }
  };

  const createPost = async () => {
    let img = "";
    if (
      document.getElementById("image").files &&
      document.getElementById("image").files[0]
    ) {
      img = await toBase64(document.getElementById("image").files[0]);
    }

    let isPrivate = document.getElementById("private").checked;
    let id = postsNr;

    let props = {
      index: id,
      card: {
        headline: document.getElementById(`new-headline-${id}`).value,
        summary: document.getElementById(`new-summary-${id}`).value,
        content: document.getElementById(`new-content-${id}`).value,
        preview_img: img,
        is_protected: isPrivate,
      },
    };

    try {
      const response = await axios.post(
        `${BASE_URL}/v1/posts/`,
        {
          headline: props.card.headline,
          summary: props.card.summary,
          content: props.card.content,
          preview_img: props.card.preview_img,
          is_protected: props.card.is_protected,
        },
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
            "Content-Type": "application/json",
            Accept: "application/json",
          },
        }
      );

      window.location.reload(true);
    } catch (errors) {
      console.error(errors);
    }
  };

  const postCard = (props) => {
    let img = "";
    if (props.card.preview_img !== "" && props.card.preview_img !== null && props.card.preview_img) {
      img = `
      <div>
      <img src="${props.card.preview_img}" />
    </div>`;
    }
    return `
  <div class="post-card" id="card-${props.index}">
    <div>
      <h2 id="headline-${props.index}">${props.card.headline}</h2>
      <div id="summary-container-${props.index}" class="post-card-summary">
        <p id="summary-${props.index}">${props.card.summary}</p>
      </div>
      <div class="post-card-footer">
      </div>
      <div>
        <input type="button" id="toggle-${props.index}" value="See more" />
      </div>
    </div>
    ${img}
  </div>
`;
  };

  const postEdit = (id) => {
    const headline = document.getElementById(`headline-${id}`).innerText;
    const summary = document.getElementById(`summary-${id}`).innerText;
    const content = document.getElementById(`content-${id}`).innerText;
    return `
  <div id="form-container-${id}">
    <form id="update-${id}" onsubmit="(event)=>{event.preventDefault();}">
      <div>
        <input type="text" id="new-headline-${id}" value="${headline}" placeholder="Headline" maxlength="80" minlength="3" required/>
      </div>
      <div>
        <input type="text" id="new-summary-${id}" value="${summary}" placeholder="Summary" maxlength="100" minlength="10" required/>
      </div>
      <div>
        <input type="text" id="new-content-${id}" value="${content}" placeholder="Content" maxlength="500" minlength="10" required/>
      </div>
      <div>
        <input type="submit" class="save" id="save-${id}" value="Save" />
        <input type="button" class="undo" id="undo-${id}" value="Undo" />
      </div>
    </form>
  </div>
  `;
  };

  const postAdd = () => {
    return `
  <div id="form-container-${postsNr}" class="form-container">
    <form id="create-${postsNr}" onsubmit="(event)=>{event.preventDefault();}">
      <div class="form-input">
        <input type="text" id="new-headline-${postsNr}" placeholder="Headline" maxlength="80" minlength="3" required/>
      </div>
      <div class="form-input">
        <input type="text" id="new-summary-${postsNr}" placeholder="Summary" maxlength="100" minlength="10" required/>
      </div>
      <div class="form-input">
        <input type="text" id="new-content-${postsNr}" placeholder="Content" maxlength="500" minlength="10" required/>
      </div>
      <div>
        <input type="file" id="image" name="image" accept="image/png, image/jpeg">
      </div>
      <div class="private-container">
        <input type="checkbox" id="private" name="private" />
        <label id="private-label" for="private">Private</label>
      </div>
      <div>
        <input type="submit" class="save" id="new-post-create" value="Post" />
        <input type="button" class="undo" id="undo-${postsNr}-initial" value="Undo" />
        <input type="button" class="undo" id="remove-image" value="Remove image" />
      </div>
    </form>
  </div>
  `;
  };

  const addPostButton = () => {
    document
      .getElementById("new-post")
      .insertAdjacentHTML(
        "afterbegin",
        ' <input type="button" id="new-post" class="new-post" value="New post" />'
      );
  };

  const showError = (message) => {
    return `
  <div class="error" id="error">
    <p>${message}</p>
  </div>
  `;
  };

  const removeError = () => {
    if (!!document.getElementById("error"))
      document.getElementById("error").remove();
  };

  const showLogin = () => {
    document.getElementById("form-register").remove();

    document
      .getElementById("container")
      .insertAdjacentHTML("afterbegin", loginForm());
  };

  const showRegister = () => {
    document.getElementById("form-login").remove();

    document
      .getElementById("container")
      .insertAdjacentHTML("afterbegin", registerForm());
  };

  const registerForm = () => {
    return `
  <div class="form-container form-register" id="form-register">
    <h2>Register</h2>
    <form id="register-form" onsubmit="(event)=>{event.preventDefault();}">
      <div class="email-container">
        <input type="email" id="email" placeholder="Email" required/>
      </div>
      <div class="password-container">
        <input type="password" id="password" placeholder="Password" required/>
      </div>
      <div class="form-footer" id="form-footer">
        <input type="submit" class="save" id="submit-register" value="Register" />
        <input type="button" class="Already have account" id="login-button" value="Already have an account?" />
      </div>
    </form>
  </div>
  `;
  };

  const loginForm = () => {
    return `
  <div class="form-container form-register" id="form-login">
    <h2>Login</h2>
    <form id="login-form" onsubmit="(event)=>{event.preventDefault();}">
      <div class="email-container">
        <input type="email" id="email" placeholder="Email" required/>
      </div>
      <div class="password-container">
        <input type="password" id="password" placeholder="Password" required/>
      </div>
      <div class="form-footer" id="form-footer">
        <input type="submit" class="save" id="submit-login" value="Login" />
        <input type="button" class="Already have account" id="register-button" value="Don't have a profile?" />
      </div>
    </form>
  </div>
  `;
  };

  const newPostButton = () => {
    return `
  <div id="new-post-container">
    <input type="button" id="new-post" class="new-post" value="New post" />
  </div>
  `;
  };

  const manageAccount = () => {
    if (
      document.getElementById("account").value === "Logout" &&
      localStorage.getItem("token")
    ) {
      localStorage.clear();
      window.location.reload(true);
      if (!!document.getElementById("form-login"))
        document.getElementById("form-login").remove();
      if (!!document.getElementById("form-register"))
        document.getElementById("form-register").remove();
    } else {
      if (!accountFlag) {
        document
          .getElementById("container")
          .insertAdjacentHTML("afterbegin", registerForm());
      } else {
        document.getElementById("form-register").remove();
      }
      accountFlag = !accountFlag;
    }
  };

  window.addEventListener("load", (event) => {
    if (localStorage.getItem("role") === "EDITOR") addPostButton();
    if (
      localStorage.getItem("role") != "EDITOR" ||
      localStorage.getItem("role") === null
    ) {
      document.getElementById("new-post").remove();
    }
    getNews();
    if (localStorage.getItem("token")) {
      document.getElementById("account").value = "Logout";
    } else {
      document.getElementById("account").value = "Account";
    }
  });

  document.addEventListener("click", function (e) {
    if (e.target.id.includes("undo")) {
      if (e.target.id.includes("initial"))
        document
          .getElementById("container")
          .insertAdjacentHTML("afterBegin", newPostButton());
      let id = e.target.id.split("-")[1];
      document.getElementById(`form-container-${id}`).remove();
    }

    if (e.target.id === "new-post") {
      document.getElementById("new-post").remove();
      document
        .getElementById("new-post-form")
        .insertAdjacentHTML("afterEnd", postAdd());
    }

    if (e.target.id === "new-post-create") {
      e.preventDefault();
      createPost();
    }

    if (e.target.id === "remove-image") {
      document.getElementById("image").value = null;
    }

    if (e.target.id === "account") {
      manageAccount();
    }

    if (e.target.id.includes("toggle")) {
      let id = e.target.id.split("-")[1];
      toggleFull(id);
    }

    if (e.target.id === "login-button") {
      showLogin();
    }

    if (e.target.id === "register-button") {
      showRegister();
    }

    if (e.target.id.includes("edit")) {
      let id = e.target.id.split("-")[1];
      document
        .getElementById(`card-${id}`)
        .insertAdjacentHTML("afterEnd", postEdit(id));
    }

    if (e.target.id === "submit-login") {
      e.preventDefault();
      removeError();
      login();
    }

    if (e.target.id === "submit-register") {
      e.preventDefault();
      removeError();
      register();
    }

    if (e.target.id.includes("save")) {
      e.preventDefault();
      let id = e.target.id.split("-")[1];
      saveForm(id);
    }
  });

  document.addEventListener("focus", function (e) {
    if (e.target.id === "email" || e.target.id === "password") {
      removeError();
    }
  });
}

export default App();
