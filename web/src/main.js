import "htmx.org";

async function onTelegramAuth(user) {
  const payload = {
    id: user.id,
    first_name: user.first_name,
    last_name: user.last_name,
    username: user.username,
    photo_url: user.photo_url,
    auth_date: user.auth_date,
    hash: user.hash,
  };

  console.log(user);

  let resp = await fetch("/api/auth/login/tg", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!resp.ok) {
    console.log(resp.json());
  } else {
    window.location.href = "/platform";
  }
}

window.onTelegramAuth = onTelegramAuth;

