#  Automatización y Control de Estacionamientos

**Proyecto de Portafolio de Título**  
**Carrera:** Analista Programador Computacional  
**Institución:** Instituto Profesional Duoc UC – Sede Concepción  
**Autores:** David Nova, Juan Espinoza, Gonzalo Honorato  
**Año:** 2025

---

## Descripción del Proyecto

Este proyecto tiene como objetivo desarrollar una solución tecnológica para **automatizar y gestionar estacionamientos** en entornos institucionales, como sedes educativas. A través del uso de tecnologías como **PWA**, **lectura de patentes (OCR)** y **códigos QR**, se busca mejorar la eficiencia, trazabilidad y experiencia de los usuarios al momento de acceder a espacios de estacionamiento.

---

## Arquitectura de Software

El sistema está basado en una **arquitectura cliente-servidor**, con separación de responsabilidades:

### Cliente (Frontend)
- **Tecnologías:** Vue 3, TypeScript, CSS
- **Funciones principales:**
  - Reserva de estacionamiento
  - Generación y escaneo de códigos QR
  - Modo offline (SQLite local)
  - Panel administrativo

### Servidor (Backend)
- **Tecnologías:** Node.js o Go (por definir)
- **Funciones principales:**
  - Autenticación con Google
  - Validación de patentes con OCR (API externa)
  - API REST para consumo del cliente
  - Almacenamiento en Google Cloud

### Base de Datos
- **Tecnología:** SQLite (modo local offline)
- **Uso:** Usuarios, reservas, registros de acceso y configuración del sistema

---

## Propuesta de Desarrollo de Software

El desarrollo se realiza utilizando **Scrum** como metodología ágil principal, complementado con prácticas de **Extreme Programming (XP)** para asegurar calidad técnica.

### Fases del Proyecto
1. Levantamiento de Requisitos
2. Diseño del Sistema (UML, Wireframes)
3. Desarrollo Técnico (Frontend + Backend)
4. Pruebas y Validación
5. Documentación y Despliegue

### Herramientas Utilizadas
- **Jira:** Gestión de tareas y backlog
- **GitHub:** Control de versiones e integración continua
- **Figma:** Prototipado UI/UX
- **Docker:** Contenedores para despliegue en VPS
- **Confluence:** Documentación técnica

---

## Plan de Gestión de Riesgos

Se identificaron y evaluaron 10 riesgos críticos y técnicos divididos en dos categorías:

- **Organizacionales del equipo**
- **Técnicos del producto**

Cada riesgo fue evaluado con base en su **probabilidad**, **impacto** y **nivel de riesgo**, y se asignaron responsables y planes de mitigación.  
Consulta la matriz y los registros detallados en el archivo: `fase_1/plan-gestion-riesgos.md`.

---

## Evidencias

### Diagramas UML
- Diagrama de casos de uso
- Diagrama de clases
- Diagrama de componentes
- Diagrama de despliegue

Ubicados en: `fase_1/evidencias/`

### Cronograma del Proyecto
El proyecto se organiza en **10 sprints**, cubriendo 4 meses de duración.

Incluye:
- Gantt general
- Tiempos estimados por actividad
- Responsables y entregables por sprint

Archivo: `fase_1/evidencias/cronograma.pdf`

---

## Documentación Complementaria

- [`fase_1/definicion-proyecto.pdf`](./fase_1/definicion_proyecto_de_software.pdf): Objetivos, contexto y relevancia académica
- [`fase_1/arquitectura-software.pdf`](./fase_1/arquitectura_de_software.pdf): Diseño técnico del sistema
- [`fase_1/propuesta-desarrollo.pdf`](./fase_1/propuesta_desarrollo_de_software.pdf): Planificación y metodología aplicada
- [`fase_1/plan-gestion-riesgos.pdf`](./fase_1/plan_de_gestion_de_riesgos.pdf): Registro y control de riesgos
- [`fase_1/evidencias/`](./fase_1/evidencias/): Resultados de sprints, avances, pruebas

---

## 🎓 Propósito Académico

Este proyecto ha sido desarrollado como parte del **portafolio de título** de la carrera *Analista Programador Computacional*, demostrando competencias en:

- Desarrollo de software seguro y de calidad
- Integración con tecnologías avanzadas (OCR, PWA, API)
- Aplicación de metodologías ágiles (Scrum + XP)
- Diseño de sistemas escalables con arquitectura limpia

---

## Contacto

- **David Nova:** da.nova@duocuc.cl  
- **Juan Espinoza:** juaa.espinoza@duocuc.cl  
- **Gonzalo Honorato:** g.honorato@duocuc.cl

---

> **Licencia y uso:** Este proyecto fue desarrollado exclusivamente con fines académicos. Cualquier uso, copia o reproducción debe contar con autorización expresa de los autores.

