
# goFinancial Planner

## Preamble

Embarking on this project represents a journey of growth and innovation for me. For several years, a trusty spreadsheet has been my steadfast companion in managing finances. It's been a dream to elevate this experience from static cells to a dynamic, interactive web application. Learning Golang, with its robust and strong typing, has been a transformative departure from my comfort zone of JavaScript's flexibility. This challenge has not only honed my technical skills but also broadened my programming horizons.

The incorporation of HTMX into this endeavor has been particularly enlightening. It has allowed me to delve into the realm of responsive web applications while maintaining a streamlined, lightweight footprint—eschewing the often cumbersome overhead of hefty SPA frameworks.

As I continue to develop and refine this application, it stands as a testament to the power of combining proven financial management strategies with cutting-edge technology. This is more than just a project; it's a milestone in my coding odyssey, embodying the spirit of continuous learning and the pursuit of turning aspirations into reality.

goFinancial Planner is a lightweight, easy-to-use financial planning web application designed to help users manage and plan their financial activities efficiently. Built with GoLang and SQLite for backend operations, and HTMX alongside Vanilla JavaScript for seamless frontend interactions, this application offers a responsive and intuitive user experience without the overhead of a heavy SPA framework or React.

## Features

- **Recurring Bills Management**: Track and manage your recurring bills with ease.
- **Financial Overview**: Get a quick overview of your financial commitments and plan accordingly.
- **Interactive UI**: Dynamic user interface with real-time updates, thanks to HTMX.
- **Database Support**: Persistent data storage with SQLite.

## Why HTMX?

HTMX allows us to harness the power of HTML extensions for dynamic updates, providing a rich user experience akin to that of a SPA, without the weight of a full JavaScript framework. This keeps our application nimble, fast, and accessible.

## Project Structure

- `controllers`: Contains handler functions for routing and business logic.
- `features`: Modular components of the application such as home, menus, etc.
- `static`: Holds static files like JavaScript, CSS, and client-side HTML.
- `database`: SQLite database files and schemas.

## Getting Started

### Prerequisites

- GoLang installed on your system.
- Node.js Installed for tooling and build steps
- Basic knowledge of HTML, JavaScript, and SQL.

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/monstercameron/gofinances.git
   ```
2. Navigate to the cloned repository:
   ```sh
   cd gofinances
   ```

3. Get Node dependencies:
   ```sh
   npm i
   ```

4. Build Templates:
   ```sh
   npm run templ
   ```

5. Build css:
    ```sh
    npm run tailwinds
    ```

5. Go install Dependencies:
    ```sh
    go mod tidy
    ```

6. Build the Go application (ensure your Go environment is set up):
   ```sh
   go build .
   ```

7. Run the application:
   ```sh
   ./gofinances
   ```

## Usage

Upon running the application, navigate to `http://localhost:3000` in your web browser to access goFinancial Planner. The UI is straightforward and self-explanatory, offering tabs for different financial categories like Recurring Bills, Debts, Assets, etc.

## Contributions

Contributions are what make the open-source community such a fantastic place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## Contact

Your Name - [@your_twitter](https://twitter.com/monstercameron)

Project Link: [https://github.com/monstercameron/gofinances](https://github.com/monstercameron/gofinances)

---

Thank you for choosing goFinancial Planner for your financial management needs.