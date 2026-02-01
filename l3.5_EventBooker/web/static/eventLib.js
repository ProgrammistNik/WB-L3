const EventAPI = {
    baseUrl: '',

    async fetchEvents() {
        try {
            const res = await fetch(`${this.baseUrl}/events`);
            const data = await res.json();
            if (!res.ok) throw data;
            return data;
        } catch (err) {
            console.error('Ошибка fetchEvents:', err);
            return { error: err.error || err.message || 'Ошибка сети' };
        }
    },

    async bookEvent(eventID, seats) {
        try {
            const res = await fetch(`${this.baseUrl}/events/${eventID}/book`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ seats })
            });
            const data = await res.json();
            if (!res.ok) throw data;
            return data;
        } catch (err) {
            console.error('Ошибка bookEvent:', err);
            return { error: err.error || err.message || 'Ошибка сети' };
        }
    },

    async createEvent({ name, date, capacity, paymentTTL }) {
        try {
            // Конвертируем дату в RFC3339 формат
            const dateObj = new Date(date);
            const rfc3339Date = dateObj.toISOString();
            
            const res = await fetch(`${this.baseUrl}/events`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ 
                    name, 
                    date: rfc3339Date, 
                    capacity, 
                    paymentTTL 
                })
            });
            const data = await res.json();
            if (!res.ok) throw data;
            return data;
        } catch (err) {
            console.error('Ошибка createEvent:', err);
            return { error: err.error || err.message || 'Ошибка сети' };
        }
    },

    async fetchBookings(eventID) {
        try {
            const res = await fetch(`${this.baseUrl}/events/${eventID}/bookings`);
            const data = await res.json();
            if (!res.ok) throw data;
            return data;
        } catch (err) {
            console.error('Ошибка fetchBookings:', err);
            return { error: err.error || err.message || 'Ошибка сети' };
        }
    },

    async confirmBooking(bookingID) {
        try {
            const res = await fetch(`${this.baseUrl}/events/${bookingID}/confirm`, { 
                method: 'POST' 
            });
            const data = await res.json();
            if (!res.ok) throw data;
            return data;
        } catch (err) {
            console.error('Ошибка confirmBooking:', err);
            return { error: err.error || err.message || 'Ошибка сети' };
        }
    },

    renderEvents(container, events, showSelectButton = true) {
        if (!events || events.error) {
            container.innerHTML = `<div class="error">Ошибка: ${events?.error || 'Нет данных'}</div>`;
            return;
        }
        if (events.length === 0) {
            container.innerHTML = '<div class="no-events">Событий нет</div>';
            return;
        }

        container.innerHTML = events.map(event => {
            const selectButton = showSelectButton ? `
                <button onclick="selectEvent(${event.id}, this)" class="select-event-btn">
                    Выбрать это мероприятие
                </button>
            ` : '';

            return `
                <div class="event-card">
                    <h3>${event.name}</h3>
                    <div class="event-info">
                        <strong>ID:</strong> ${event.id} | 
                        <strong>Дата:</strong> ${new Date(event.date).toLocaleString()} | 
                        <strong>Места:</strong> ${event.freeSeats} / ${event.capacity}
                    </div>
                    <div class="${event.freeSeats === 0 ? 'status-expired' : 'status-pending'}">
                        ${event.freeSeats === 0 ? 'Мест нет' : 'Есть свободные места'}
                    </div>
                    ${selectButton}
                </div>
            `;
        }).join('');
    },

    renderBookings(container, bookings, eventId = null, showConfirmButton = true) {
        container.innerHTML = '';

        if (!bookings || bookings.error) {
            container.innerHTML = `<div class="error">Ошибка: ${bookings?.error || 'Нет данных'}</div>`;
            return;
        }
        if (bookings.length === 0) {
            container.innerHTML = '<div class="no-bookings">Бронирований нет</div>';
            return;
        }

        const now = new Date();

        container.innerHTML = bookings.map(booking => {
            const expiresAt = new Date(booking.expiresAt);
            const isExpired = !booking.paid && expiresAt < now;

            let statusClass = '';
            let statusText = '';

            if (booking.paid) {
                statusClass = 'status-paid';
                statusText = '✅ Оплачено';
            } else if (isExpired) {
                statusClass = 'status-expired';
                statusText = '❌ Просрочено';
            } else {
                statusClass = 'status-pending';
                statusText = '⏳ Ожидает оплаты';
            }

            let buttonHtml = '';
            if (showConfirmButton && !booking.paid && !isExpired) {
                buttonHtml = `
                    <button class="confirm-btn" onclick="window.confirmBooking(${booking.id}, ${eventId})">
                        Подтвердить оплату
                    </button>
                `;
            }

            return `
                <div class="booking-card ${booking.paid ? 'confirmed-now' : ''}" 
                     style="border: 2px solid ${booking.paid ? '#28a745' : isExpired ? '#dc3545' : '#ffc107'};">
                    <div class="booking-info">
                        <strong>Бронь #${booking.id}</strong> | Мест: ${booking.seats} | 
                        <span class="${statusClass}">${statusText}</span>
                    </div>
                    <div class="booking-dates">
                        Создано: ${new Date(booking.createdAt).toLocaleString()} | 
                        Истекает: ${expiresAt.toLocaleString()}
                    </div>
                    ${buttonHtml}
                </div>
            `;
        }).join('');
    },

    startAutoRefresh(container, eventId, interval = 10000, showConfirmButton = true) {
        return setInterval(async () => {
            console.log('Автообновление бронирований...');
            const bookings = await this.fetchBookings(eventId);
            this.renderBookings(container, bookings, eventId, showConfirmButton);
        }, interval);
    }
};

// Глобальная функция для кнопки подтверждения (только для пользователя)
window.confirmBooking = async function(bookingId, eventId) {
    if (!confirm('Подтвердить оплату брони #' + bookingId + '?')) {
        return;
    }

    const result = await EventAPI.confirmBooking(bookingId);
    if (result.error) {
        alert('Ошибка: ' + result.error);
    } else {
        // Обновляем бронирования на ВСЕХ страницах где они отображаются
        updateAllBookingDisplays(eventId);
        
        // Показываем сообщение об успехе
        showSuccessMessage('✅ Бронь успешно подтверждена!');
    }
};

// Функция для обновления всех отображений бронирований
async function updateAllBookingDisplays(eventId) {
    // Обновляем на странице пользователя
    const userContainer = document.getElementById('my-bookings');
    if (userContainer) {
        const bookings = await EventAPI.fetchBookings(eventId);
        EventAPI.renderBookings(userContainer, bookings, eventId, true);
    }
    
    // Обновляем на странице админа
    const adminContainer = document.getElementById('bookings-list');
    if (adminContainer) {
        const bookings = await EventAPI.fetchBookings(eventId);
        EventAPI.renderBookings(adminContainer, bookings, eventId, false);
    }
}

// Функция для показа сообщения об успехе
function showSuccessMessage(message) {
    // Показываем на странице пользователя
    const userContainer = document.getElementById('my-bookings');
    if (userContainer) {
        const successMsg = document.createElement('div');
        successMsg.className = 'confirmation-success';
        successMsg.textContent = message;
        userContainer.prepend(successMsg);
        setTimeout(() => successMsg.remove(), 3000);
    }
    
    // Показываем на странице админа
    const adminContainer = document.getElementById('bookings-list');
    if (adminContainer) {
        const successMsg = document.createElement('div');
        successMsg.className = 'confirmation-success';
        successMsg.textContent = message;
        adminContainer.prepend(successMsg);
        setTimeout(() => successMsg.remove(), 3000);
    }
}

// Функция выбора мероприятия
window.selectEvent = function(eventId, element) {
    document.getElementById('booking-event-id').value = eventId;
    document.getElementById('my-bookings-event-id').value = eventId;
    
    // Подсвечиваем выбранное мероприятие
    document.querySelectorAll('.event-card').forEach(card => {
        card.style.border = '1px solid #ccc';
    });
    
    // Используем element для поиска родительской карточки
    if (element && element.closest) {
        element.closest('.event-card').style.border = '2px solid #007bff';
    }
    
    // Загружаем бронирования для выбранного мероприятия
    const container = document.getElementById('my-bookings');
    if (container) {
        EventAPI.fetchBookings(eventId).then(bookings => 
            EventAPI.renderBookings(container, bookings, eventId, true)
        );
    }
};