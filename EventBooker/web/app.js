const API_BASE = "";

const pageType = document.body.dataset.page;

const eventsList = document.getElementById("events-list");
const eventDetails = document.getElementById("event-details");
const messageBox = document.getElementById("message");
const createEventForm = document.getElementById("create-event-form");
const loadEventsBtn = document.getElementById("load-events-btn");
const confirmBookingForm = document.getElementById("confirm-booking-form");

let selectedEventId = null;

function showMessage(text, type = "success") {
    if (!messageBox) return;
    messageBox.textContent = text;
    messageBox.className = `message show ${type}`;
    setTimeout(() => {
        messageBox.className = "message";
    }, 3000);
}

async function request(url, options = {}) {
    const res = await fetch(url, {
        headers: {
            "Content-Type": "application/json",
        },
        ...options,
    });

    if (!res.ok) {
        let errorText = `HTTP ${res.status}`;
        try {
            const data = await res.json();
            errorText = data.message || JSON.stringify(data);
        } catch {
            errorText = await res.text();
        }
        throw new Error(errorText || "Request failed");
    }

    if (res.status === 204) return null;
    return res.json();
}

async function loadEvents() {
    if (!eventsList) return;

    try {
        const events = await request(`${API_BASE}/events`);
        renderEvents(events || []);
    } catch (err) {
        showMessage(`Ошибка загрузки событий: ${err.message}`, "error");
    }
}

function renderEvents(events) {
    if (!eventsList) return;

    if (!events || events.length === 0) {
        eventsList.innerHTML = "<div>Событий пока нет</div>";
        return;
    }

    eventsList.innerHTML = events
        .map((event) => {
            const isAdmin = pageType === "admin";

            return `
        <div class="event-item">
          <h3>${event.title}</h3>
          <div class="event-meta">ID: ${event.id}</div>
          <div class="event-meta">Описание: ${event.description}</div>
          <div class="event-meta">Дата: ${new Date(event.date).toLocaleString()}</div>
          <div class="event-meta">Свободных мест: ${event.availableSeats}/${event.totalSeats}</div>
          <div class="event-meta">Цена: ${event.price}</div>
          <div class="event-meta">Оплата: ${event.paymentRequired ? "Да" : "Нет"}</div>

          <div class="event-actions">
            <button onclick="loadEventDetails('${event.id}')">Открыть детали</button>
            ${
                isAdmin
                    ? ""
                    : `<button onclick="bookEvent('${event.id}')">Забронировать</button>`
            }
          </div>
        </div>
      `;
        })
        .join("");
}

async function loadEventDetails(eventId) {
    selectedEventId = eventId;

    try {
        const data = await request(`${API_BASE}/events/${eventId}`);
        renderEventDetails(data);
    } catch (err) {
        showMessage(`Ошибка загрузки деталей: ${err.message}`, "error");
    }
}

function renderEventDetails(data) {
    if (!eventDetails) return;

    const event = data.event || data.Event;
    const bookings = data.bookings || data.Bookings || [];
    const freeSeats = data.freeSeats ?? data.FreeSeats;
    const totalBooked = data.totalBooked ?? data.TotalBooked;

    eventDetails.innerHTML = `
    <div class="details-block">
      <div>
        <h3>${event.title}</h3>
        <div>ID: ${event.id}</div>
        <div>Описание: ${event.description}</div>
        <div>Дата: ${new Date(event.date).toLocaleString()}</div>
        <div>Свободных мест: ${freeSeats}</div>
        <div>Активных броней: ${totalBooked}</div>
      </div>

      <div>
        <h3>Брони</h3>
        <div class="booking-list">
          ${
        bookings.length === 0
            ? "<div>Броней пока нет</div>"
            : bookings
                .map(
                    (b) => `
                  <div class="booking-item">
                    <div><strong>ID:</strong> ${b.id}</div>
                    <div><strong>User:</strong> ${b.username}</div>
                    <div><strong>Status:</strong> ${b.status}</div>
                    <div><strong>Created:</strong> ${new Date(b.createdAt).toLocaleString()}</div>
                    <div><strong>Expires:</strong> ${b.expiredAt ? new Date(b.expiredAt).toLocaleString() : "-"}</div>
                    <div><strong>Confirmed at:</strong> ${b.confirmed_at ? new Date(b.confirmed_at).toLocaleString() : "-"}</div>
                  </div>
                `
                )
                .join("")
    }
        </div>
      </div>
    </div>
  `;
}

async function bookEvent(eventId) {
    const username = prompt("Введите username для бронирования:");
    if (!username) return;

    try {
        const booking = await request(`${API_BASE}/events/${eventId}/book`, {
            method: "POST",
            body: JSON.stringify({ username }),
        });

        showMessage(`Бронь создана: ${booking.id}`);
        await loadEvents();
        await loadEventDetails(eventId);
    } catch (err) {
        showMessage(`Ошибка бронирования: ${err.message}`, "error");
    }
}

if (createEventForm) {
    createEventForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const payload = {
            title: document.getElementById("title").value,
            description: document.getElementById("description").value,
            date: new Date(document.getElementById("date").value).toISOString(),
            totalSeats: Number(document.getElementById("totalSeats").value),
            availableSeats: Number(document.getElementById("totalSeats").value),
            price: Number(document.getElementById("price").value),
            paymentRequired: document.getElementById("paymentRequired").checked,
        };

        try {
            await request(`${API_BASE}/events`, {
                method: "POST",
                body: JSON.stringify(payload),
            });

            createEventForm.reset();
            showMessage("Событие создано");
            await loadEvents();
        } catch (err) {
            showMessage(`Ошибка создания события: ${err.message}`, "error");
        }
    });
}

if (confirmBookingForm) {
    confirmBookingForm.addEventListener("submit", async (e) => {
        e.preventDefault();

        const bookingId = document.getElementById("confirmBookingId").value.trim();
        if (!bookingId) return;

        try {
            await request(`${API_BASE}/events/${bookingId}/confirm`, {
                method: "POST",
            });

            showMessage("Бронь подтверждена");

            if (selectedEventId) {
                await loadEventDetails(selectedEventId);
                await loadEvents();
            }
        } catch (err) {
            showMessage(`Ошибка подтверждения: ${err.message}`, "error");
        }
    });
}

if (loadEventsBtn) {
    loadEventsBtn.addEventListener("click", loadEvents);
}

loadEvents();