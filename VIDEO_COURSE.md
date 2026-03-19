# Video Course: DocProcessor Deep Dive

5-episode series covering DocProcessor from fundamentals to advanced usage.

---

## Episode 1: Introduction & Architecture (15 min)

### Script Outline

1. **Opening** (2 min)
   - What is DocProcessor and why it exists
   - Role in the HelixQA ecosystem
   - Problem statement: manual QA doesn't scale

2. **Architecture Overview** (5 min)
   - Package layout: loader, feature, coverage, docgraph, llm, config
   - Processing pipeline visualization
   - Thread safety design with RWMutex

3. **Quick Start Demo** (5 min)
   - Clone and build
   - Run CLI on sample docs
   - Show feature map output

4. **Key Design Decisions** (3 min)
   - LLMAgent injection (no hard dependencies)
   - Deterministic feature IDs
   - Heuristic fallback when no LLM

---

## Episode 2: Document Loading & Parsing (15 min)

### Script Outline

1. **Loader Interface** (3 min)
   - `LoadDir`, `LoadFile`, `SupportedFormats`
   - Context cancellation support
   - MaxFileSize protection

2. **Markdown Parser** (4 min)
   - Heading extraction with regex
   - Section content segmentation
   - Link extraction and deduplication

3. **YAML Parser** (3 min)
   - Generic unmarshalling
   - Title extraction from `title` or `name` keys
   - Section creation from top-level keys

4. **Scanner Implementation** (3 min)
   - Recursive directory walking
   - Hidden directory skipping
   - Extension normalization (yml -> yaml)

5. **Security** (2 min)
   - Path traversal protection
   - File size limits
   - Malformed input handling

---

## Episode 3: Feature Extraction & FeatureMap (15 min)

### Script Outline

1. **Feature Model** (3 min)
   - Feature struct: ID, Name, Category, Platforms, Priority
   - FeatureCategory constants
   - Deterministic ID generation

2. **Heuristic Extraction** (3 min)
   - Section-based extraction
   - Content length thresholds
   - Default platform assignment

3. **LLM-Powered Extraction** (4 min)
   - LLMAgent interface design
   - Prompt templates with versioning
   - ExtractFeatures, ClassifyFeature flow

4. **FeatureMap Building** (3 min)
   - AddFeature with category/platform indexing
   - DocGraph integration
   - FeatureByID, FeaturesForPlatform queries

5. **Merge & Deduplication** (2 min)
   - Merging multiple feature maps
   - ID-based deduplication
   - DocGraph merging

---

## Episode 4: Coverage Tracking (15 min)

### Script Outline

1. **CoverageTracker Interface** (3 min)
   - RegisterFeature, MarkVerified, MarkFailed, MarkSkipped
   - Coverage reports by platform and category
   - Export for serialization

2. **Thread-Safe Implementation** (4 min)
   - sync.RWMutex usage
   - Read path: RLock for Coverage(), CoverageByPlatform(), Unverified()
   - Write path: Lock for MarkVerified(), MarkFailed(), MarkSkipped()

3. **Verification States** (3 min)
   - State machine: unverified -> verified/failed/skipped
   - State transitions and overwrites
   - Evidence and Issue tracking

4. **Reporting** (3 min)
   - CoverageReport aggregation
   - Platform and category breakdowns
   - Issue collection

5. **Stress Testing** (2 min)
   - Concurrent verification from multiple platform workers
   - Race condition detection
   - High-contention scenarios

---

## Episode 5: DocGraph & Integration (15 min)

### Script Outline

1. **DocGraph** (4 min)
   - Directed graph of document links
   - Node and edge management
   - Neighbor and incoming edge queries

2. **Export Formats** (3 min)
   - JSON export and import (round-trip)
   - Mermaid diagram generation
   - ID sanitization for Mermaid

3. **Full Pipeline Integration** (4 min)
   - Load -> Parse -> Extract -> Build -> Enrich -> Track
   - Mock LLMAgent for testing
   - E2E test walkthrough

4. **Configuration** (2 min)
   - .env file loading
   - Default values
   - Environment variable mapping

5. **What's Next** (2 min)
   - HelixQA integration via bridge adapters
   - LLMOrchestrator Agent -> LLMAgent adapter
   - Autonomous QA session pipeline
