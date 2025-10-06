# 📅 Sentiric Scheduling Service - Mantık ve Akış Mimarisi

**Stratejik Rol:** Kullanıcılar veya Agent'lar adına takvimde (Google Calendar, Outlook) randevu oluşturma ve değiştirme işlemlerini yöneten merkezi soyutlama katmanı.

---

## 1. Temel Akış: Randevu Oluşturma (CreateAppointment)

Bu servis, gelen randevu isteğini platformdan bağımsız (agnostik) bir şekilde alır ve yapılandırılmış takvim sağlayıcısına (Adapter) yönlendirir.

```mermaid
graph TD
    A[Agent Service] -- gRPC: CreateAppointment(...) --> B(Scheduling Service)
    
    Note over B: 1. Adaptör Seçimi (Google/Outlook Adapter)
    B --> C{Google Calendar Adaptörü};
    C -- OAuth2 API Çağrısı --> Google[Harici Google API];
    Google -- Confirmation --> C;
    
    Note over C: Randevu ID'si Alınır.
    C --> B;
    B -- Response --> A;
```

## 2. Adaptör Mimarisi

Scheduling Service, CalendarAdapter yapılandırmasına göre uygun adaptörü seçer.

* Desteklenecekler: Google Calendar, Microsoft Outlook/Exchange
* API Gereksinimi: Çoğu takvim entegrasyonu OAuth2 veya API Anahtarı gerektirir, bu kimlik bilgileri config'de yönetilmelidir.