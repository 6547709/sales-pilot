# UI Optimizations Design (2026-03-29)

## Overview
Optimize the sales-pilot application based on user feedback regarding login timeouts and mobile navigation layout.

## 1. Login Timeout Extension
**Problem:** The current login session is too short (24 hours).
**Solution:** Modify the backend JWT mechanism to extend the token TTL to 14 days (2 weeks).
**Implementation Details:**
- Update `TokenTTL` constant in `backend/internal/auth/jwt.go` from `24 * time.Hour` to `14 * 24 * time.Hour`.

## 2. Mobile Navigation Uncollapse
**Problem:** Mobile layout hides navigation links ("首页", "方案库", "退出登录") inside a Sheet (hamburger menu), adding unnecessary friction since there is ample space in the top row.
**Solution:** Eliminate the `<Sheet>` component from the marketing shell for public links and basic auth actions. Instead, display the links inline dynamically scaling them or placing them side-by-side using Flexbox layout, ensuring they fit within the screen width on typical mobile devices.

## 3. Smart Sticky Navbar
**Problem:** The sticky top navigation occupies valuable screen real estate when users are reading or scrolling down.
**Solution:** Implement "Smart Sticky" (Option A) behavior using React hooks (`useEffect` on scroll event or custom hook like `useScrollDirection`).
- When the user scrolls down, the header translates upward to disappear (`transform: translateY(-100%)`).
- When the user scrolls up, the header translates back down to become visible (`transform: translateY(0)`).
- Transition properties will ensure smooth sliding.
