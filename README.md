# AWASM Portfolio

![Coverage Badge](https://raw.githubusercontent.com/trianalab/awasm-portfolio/badges/.badges/main/coverage.svg) ![Latest Release](https://img.shields.io/github/v/release/trianalab/awasm-portfolio) ![License](https://img.shields.io/github/license/trianalab/awasm-portfolio)  

AWASM Portfolio is a WebAssembly-powered application that emulates a console, enabling users to interact with the developer's resume using Kubernetes-like commands. This innovative client-side architecture combines frontend and backend functionality for seamless and dynamic user interactions.

![Demo](images/demo.gif)

But you can also check the live demo [here](https://trianalab.github.io/awasm-portfolio/).

## Deployment
### Single-Command Deployments

1. **Locally**: Use the provided Makefile:
   ```bash
   make run
   ```
   Access the application at `http://127.0.0.1:8000`.

2. **Docker**: Pull and run the latest release image:
   ```bash
   docker run -p 8000:80 ghcr.io/trianalab/trianalab/awasm-portfolio:$(curl -s https://api.github.com/repos/trianalab/awasm-portfolio/releases/latest | jq -r .tag_name)
   ```
   Access the application at `http://127.0.0.1:8000`.

## Architecture Overview

![Architecture](images/architecture.png)

### Frontend
The frontend features two primary modes of interaction, both synchronized with the WebAssembly backend:

1. **CLI Mode**  
   - Emulates a web terminal using xterm.js.  
   - Supports `kubectl` commands to interact with a virtual cluster.

2. **UI Mode**  
   - Displays data as interactive, physics-based bubbles.  
   - Offers a dynamic, modern visualization of CV information.

### Backend
The backend, written in Go and compiled to WebAssembly, runs entirely within the browser. It processes commands and manages data using a layered architecture:

- **Cmd Package**: Handles command parsing (via Cobra) and validation.  
- **Services Package**: Implements business logic and data processing.  
- **Repository Package**: Manages an in-memory key-value store.  
- **Factory Package**: Creates mock resources with randomized data (using gofakeit).  
- **Models Package**: Defines Kubernetes-like CRD models for CV elements.  
- **Preload Package**: Loads initial CV data for immediate use.  
- **UI Package**: Serializes data into YAML, JSON, or tables for consistent display.  
- **Util Package**: Ensures data normalization and utility operations.  

## Key Features
- **Integrated Architecture**: Entirely client-side for high performance and offline capability.  
- **Dual Interaction Modes**: Switch between CLI and UI modes with real-time synchronization.  
- **Kubernetes-Inspired Commands**: Use familiar `kubectl` syntax to explore the CV.  
- **Dynamic Visuals**: Physics-based UI for engaging data visualization.  
- **Extensibility**: Easily customizable models and data generation.  

## Contributing
Contributions are welcome! Open an issue or submit a pull request to share your ideas or fixes.
