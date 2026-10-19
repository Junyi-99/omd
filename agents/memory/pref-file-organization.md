# File Organization Preferences and Habits

## Overview
This document records file organization preferences and habits for future reference. When asked to organize files, refer to this document to understand the preferred structure and categorization approach.

## Chosen Organization System: Plan 3 - Hybrid Organization

The user prefers a **hybrid organization approach** that separates:
- **Projects** (code/development work) from **Media** (non-code files)
- Clear domain-based categorization within Projects
- Type-based organization within Media

## Directory Structure

### Main Categories

```
Desktop/
├── Projects/          # All code projects, organized by domain
├── Media/             # All non-code files (photos, docs, videos, etc.)
├── Academic/          # Coursework, homework, research materials
└── Archives/          # Test folders, old projects, inactive code
```

## Detailed Organization Rules

### 1. Projects/ Structure

Projects are organized by **domain/purpose**, not by file type:

#### Projects/Research/
- Academic research projects
- Paper-related code (paperdebugger-* variants)
- AI/ML research (CoT-with-Supervision, academic-paper-mcp-server)
- Conference submissions (neurips25)
- Research tools and utilities

**Examples:**
- `academic-paper-mcp-server`
- `CoT-with-Supervision`
- `paperdebugger-project`, `paperdebugger-analysis`, `paperdebugger-*` (all variants)
- `neurips25`
- `research`

#### Projects/Trading/
- Trading bots and strategies
- Quantitative analysis projects
- Market-related code
- Polymarket projects
- Gambling/betting related code

**Examples:**
- `trading_backup`, `trading-bot-*`, `trading-coefficient`
- `quant`, `quant2`, `quant3`
- `gambling`
- `polymarket-bot`, `polymarket-archive-*`
- `bid_ask`
- `凯利公式实验` (Kelly formula experiments)

#### Projects/Web/
- Web development projects
- Personal websites and blogs
- Frontend projects
- Web applications

**Examples:**
- `hou-junyi.github.io`, `junyi-blog`
- `website`, `paperdebugger-website`
- `nextjs-boilerplate`
- `foodcourt-src`, `hack`
- `Sub-Store`, `Sub-Store-Front-End`
- `rateyourpaper.ai`, `xacc_web`, `icpp2026.github.io`

#### Projects/Infrastructure/
- DevOps and infrastructure projects
- Kubernetes/cluster management
- Server configuration and setup
- Network tools

**Examples:**
- `k3s_rpi`, `my-k8s-cluster`
- `servers`
- `xtra-*` (all variants: xtra-playbooks, xtra-server-*, xtramcp)
- `Shadowrocket-Modules`

#### Projects/Tools/
- Development tools and utilities
- SDKs and libraries
- API-related projects
- Configuration management

**Examples:**
- `ChatGPT-API-Scanner`
- `openai-sdk-streaming`
- `vscode-copilot-chat`
- `python-sdk`
- `oh-my-dotfiles`

#### Projects/ (root level for other projects)
- Large frameworks/libraries (MNN, hugo, monorepo)
- Miscellaneous projects that don't fit other categories
- Individual code files

**Examples:**
- `bean`, `scaffold`, `monorepo`
- `MNN`, `hugo`, `hugoThemesSiteBuilder`
- `mlx-swift-examples`
- `on-chain`
- Individual `.py`, `.ts`, `.js`, `.go`, `.cpp` files

### 2. Media/ Structure

Media files are organized by **type**, with photos having additional sub-organization:

#### Media/Photos/
Photos are further organized by **context/date**:

- **Media/Photos/2025/** - Photos from 2025
  - `2025 新加坡国庆节` (Singapore National Day 2025)

- **Media/Photos/2026/** - Photos from 2026
  - `2026`

- **Media/Photos/Travel/** - Travel-related photos
  - `厦门` (Xiamen)
  - `实里达` (Seletar)
  - `新西兰` (New Zealand)
  - `PasirRis`

- **Media/Photos/Personal/** - Personal photos
  - `myself`
  - `翠鸟` (Kingfisher)
  - `Ridge Computing`, `RidgeComputing`
  - Individual photos (IMG_*, photo*, Selfie*, 微信图片*)

- **Media/Photos/Camera/** - Camera RAW files and camera folders
  - `100MSDCF` (Sony camera folder)
  - `DCIM` (Digital Camera Images)
  - `35mm夜景` (35mm night scenes)
  - Individual `.ARW`, `.HEIC` files

- **Media/Photos/Untitled Export/** - Exported photo collections
  - `Untitled Export`

#### Media/Screenshots/
- All screenshot files (Screenshot*.png, SCR-*.png, SCR-*.jpeg)
- Simulator screenshots (Simulator Screenshot*.png)
- Screen recordings moved to Videos instead

#### Media/Documents/
Documents organized by **topic/purpose**:

- **Media/Documents/Academic/** - Academic papers, research documents
  - Research papers (PAPER.pdf, how-ai-trust-human.pdf, etc.)
  - Academic presentations
  - Course-related PDFs
  - Diagrams and charts (non-photo images)

- **Media/Documents/Personal/** - Personal documents
  - ID cards, passports
  - Bank statements (CardStatement_*.pdf)
  - Travel itineraries (ITINERARY.pdf, NZ_Travel_Itinerary.docx)
  - Visa applications
  - Personal letters and documents

- **Media/Documents/Legal/** - Legal and official documents
  - Constitutions, terms and conditions
  - Official announcements (公報-*)
  - Legal paperwork

#### Media/Presentations/
- All PowerPoint files (`.pptx`)
- Academic presentations
- Work presentations
- Seminar slides

#### Media/Videos/
- All video files (`.mov`, `.mp4`)
- Screen recordings
- Demo videos
- Personal videos

### 3. Academic/ Structure

Academic folder contains:
- Coursework folders (`cs6216_*`, `hw4`)
- Lab work (`Tokka-Lab-*`)
- LaTeX projects (`latex`, `latex-ai-trust-human`, `latex-projects`)
- Research materials (`Markdown 资料`, `文章收藏` - article collection)
- Jupyter notebooks (`.ipynb`)
- Text files (`.txt`)
- Course-related files

### 4. Archives/ Structure

Archives folder contains:
- Test folders (`test`, `test2`, `test3`, `test4`, `test-faro`)
- Old/inactive projects
- `node_modules` folders
- Experimental folders (`exp.oct.9`)
- Temporary folders (`section`, `blog_nuc`, `typora`, `chatlog`)
- Large archived projects (`zhang_rong`)
- Downloaded web content folders
- Other inactive or experimental code

## File Categorization Rules

### Individual Files at Root Level

When organizing individual files found at the desktop root:

1. **Screenshots** → `Media/Screenshots/`
   - Pattern: `Screenshot*.png`, `SCR-*.png`, `SCR-*.jpeg`, `Simulator Screenshot*.png`

2. **Photos** → `Media/Photos/Personal/` or `Media/Photos/Camera/`
   - Personal photos: `IMG_*`, `photo*`, `Selfie*`, `微信图片*`
   - RAW files: `.ARW`, `.HEIC` → `Media/Photos/Camera/`

3. **Videos** → `Media/Videos/`
   - All `.mov`, `.mp4` files

4. **Presentations** → `Media/Presentations/`
   - All `.pptx` files

5. **Documents** → `Media/Documents/` (subfolder by topic)
   - Academic papers → `Media/Documents/Academic/`
   - Personal docs → `Media/Documents/Personal/`
   - Legal docs → `Media/Documents/Legal/`

6. **Code Files** → `Projects/`
   - Individual `.py`, `.ts`, `.js`, `.go`, `.cpp`, `.hpp` files → `Projects/` root
   - Config files (`.yaml`, `.yml`) → `Projects/Infrastructure/`
   - Notebooks (`.ipynb`) → `Academic/`

7. **Data Files** → `Projects/`
   - `.csv`, `.json` files → `Projects/` root

8. **Archives** → `Archives/`
   - `.zip`, `.tar.gz` files

## User Habits and Preferences Observed

### 1. Language Usage
- Comfortable with **Chinese folder names** (新加坡国庆节, 厦门, 实里达, 新西兰, 翠鸟, 文章收藏, 凯利公式实验)
- Mixes English and Chinese in folder/file names
- **Preference**: Preserve original folder names, including Chinese characters

### 2. Project Organization
- Creates **multiple variants** of projects (e.g., `paperdebugger-*` has many variants: project, analysis, assets, blog, l0-l1, mcp, website, copy)
- Keeps **test folders** separate (test, test2, test3, test4)
- Has **backup folders** (trading_backup, paperdebugger copy)
- **Preference**: Keep all variants together in the same category, don't merge them

### 3. Photo Organization
- Organizes photos by **date** (2025, 2026 folders)
- Organizes photos by **location** (travel folders with Chinese place names)
- Keeps **camera RAW files** separate (100MSDCF, DCIM)
- Has **personal photo collections** (myself, 翠鸟)
- **Preference**: Maintain date/location context for photos

### 4. Academic Work
- Separates **coursework** from **research**
- Keeps **LaTeX projects** separate
- Has **article collections** (文章收藏)
- **Preference**: Academic folder for coursework, Projects/Research for research code

### 5. File Naming
- Uses **descriptive names** with dates (Screenshot 2025-09-02, SCR-20251218)
- Uses **project-specific prefixes** (paperdebugger-*, trading-*, quant*)
- Uses **abbreviations** (CoT, MCP, SDK)
- **Preference**: Preserve original naming conventions

### 6. Development Practices
- Works with **multiple languages** (Python, TypeScript, Go, C++, Swift)
- Has **infrastructure projects** (k3s, k8s, servers)
- Uses **version control** (github.io projects suggest Git usage)
- **Preference**: Organize by domain/purpose, not by language

### 7. Archive Behavior
- Keeps **test folders** for experimentation
- Has **node_modules** on desktop (should be archived)
- Keeps **old projects** (zhang_rong, Digital Science folders)
- **Preference**: Archive inactive/test projects, but keep them accessible

## Implementation Approach

### Script Safety Features
- **No overwriting**: Check if destination exists before moving
- **Preserve structure**: Move folders as-is, don't modify contents
- **Error handling**: Skip missing files gracefully
- **Idempotent**: Safe to run multiple times
- **Verbose output**: Show what's being moved

### Key Principles
1. **Preserve original names**: Never rename folders or files
2. **Maintain context**: Keep related files together when possible
3. **Respect existing structure**: Don't modify folder contents
4. **Clear separation**: Projects vs Media vs Academic vs Archives
5. **Domain-based projects**: Organize by purpose, not file type
6. **Type-based media**: Organize media by file type with context

## Future Organization Requests

When asked to organize files in the future:

1. **Reference this document** for structure and preferences
2. **Ask for confirmation** if structure differs significantly
3. **Preserve original names** and folder structures
4. **Create scripts** that are safe and reviewable
5. **Respect Chinese folder names** and mixed-language usage
6. **Keep project variants together** in same category
7. **Maintain date/location context** for photos
8. **Separate active projects** from archives

## Notes

- User prefers to **review scripts** before execution
- User wants **manual control** over file organization
- User values **preservation** of existing structures
- User uses **both English and Chinese** in file organization
- User has **many active projects** across different domains
- User keeps **detailed photo collections** with context
