# 项目面试亮点

## 用户模块

在用户模块的设计中，我采取了以下关键步骤来优化安全性、性能和用户体验：

1. **身份校验与令牌管理**：
   - **短Token（身份校验Token）**：
      - 用于用户身份验证，每次用户请求时在HTTP头中携带。
      - 具有较短的有效期，如5分钟，以增强安全性。
      - 存储在客户端，如浏览器Cookie。
   - **长Token（刷新Token）**：
      - 用于刷新短Token，确保用户持续登录状态。
      - 具有较长的有效期，如一周，减少用户重新登录的频率。
      - 存储在服务端，通常通过Redis等缓存服务实现。
2. **令牌生成与验证**：
   - 使用JWT（JSON Web Tokens）技术生成和验证Token。
   - 短Token和长Token都包含用户ID、SSID（会话ID）等信息。
   - 短Token包含用户代理和内容类型等信息，用于追踪用户行为。
3. **安全性与性能**：
   - 短Token的有效期短，即使被泄露，威胁窗口也较短。
   - 长Token在服务端存储，即使泄露，也不会直接暴露用户身份。
   - 减少了对数据库的直接访问，提高了性能。
4. **用户界面**：
   - 设计直观的用户登录和注册界面，简化用户操作。
   - 实现自动登录功能，使用户在刷新页面或重新加载时无需重新登录。
5. **数据同步**：
   - 实现Token的有效期管理，自动刷新短Token。
   - 确保长Token在用户会话结束时失效，以防止滥用。
6. **异常处理**：
   - 实现Token过期或无效时的错误处理，返回清晰的错误信息。
   - 监控Token的生成和使用情况，及时发现和处理异常。

## 帖子模块

在帖子模块的设计中，我采取了以下关键步骤来优化性能和用户体验：

1. **数据库选择**：
   - 线上库（`Posts`）：使用MongoDB存储已发布的帖子，利用其灵活的文档模型和高效的读写性能
   - 制作库（`Drafts`）：使用MySQL存储未发布的帖子，利用其支持事务和严格数据一致性的特点
2. **数据同步**：
   - 当用户发布帖子时，系统将帖子从制作库同步到线上库，确保帖子的状态和内容准确无误
3. **查询逻辑**：
   - 已发布的帖子通过线上库查询，提供快速、高效的访问
   - 未发布的帖子通过制作库查询，仅允许帖子的所有者访问
4. **用户界面**：
   - 设计用户友好的界面，包括“保存草稿”和“发布”按钮，便于用户管理帖子
5. **安全性**：
   - 实现严格的权限控制，确保只有帖子的所有者可以访问和编辑未发布的帖子
6. **性能优化**：
   - 考虑使用缓存机制来提高查询效率，同时确保同步操作不会影响系统性能

### 为什么选择分开存储(MongoDB/MySQL)

1. **技术优势**:
   - **MongoDB**：
      - 适合存储和查询文档型数据，如帖子的内容和评论
      - 提供灵活的查询能力，支持复杂的聚合操作
      - 适用于读写操作频繁的场景
   - **MySQL**：
      - 适合存储结构化数据，如用户信息和帖子元数据
      - 支持事务处理，确保数据一致性
      - 提供丰富的数据管理和备份工具
2. **安全性**：
   - 可以根据不同的安全需求，对两个数据库实施不同的访问控制策略

## 榜单模块

在榜单模块的设计中，我采取了以下关键步骤来优化性能和用户体验：

1. **分批处理**：
   - **批处理大小**：默认设置每次分页处理的帖子数量为100，保证在处理大量数据时不会占用过多资源。
   - **排名数量**：默认设置要计算并返回的排名前100帖子的数量，确保榜单的结果能够快速响应并呈现给用户。
2. **计算逻辑**：
   - **帖子评分函数**：根据点赞数和更新时间计算每个帖子的分数，确保评分逻辑能够准确反映帖子的受欢迎程度和时效性。
   - **优先队列**：使用优先队列对帖子进行排序，确保计算效率和内存使用的优化。
   - **时间截止**：将七天前的时间作为截止时间，确保榜单数据的时效性和动态更新。
3. **数据获取**：
   - **分页获取帖子**：通过分页方式从存储库中获取已发布的帖子，保证在大数据量时的获取效率。
   - **获取交互数据**：根据帖子ID列表获取交互数据（如点赞数），确保数据的完整性和准确性。
4. **安全性与性能**：
   - **日志记录**：使用zap日志库记录关键操作和错误，便于调试和监控。
   - **异常处理**：处理在获取帖子和交互数据时可能出现的错误，确保系统的健壮性。

### 技术优势

1. **技术选型**：
   - **优先队列**：使用优先队列进行帖子排序，确保在处理大量数据时仍能保持高效的排序和检索性能。
   - **分页处理**：分页获取和处理帖子，避免一次性加载大量数据导致的内存溢出和性能问题。
   - **分数计算函数**：采用灵活的分数计算函数，可以根据业务需求进行调整，确保排名结果的准确性和合理性。
2. **安全性**：
   - **日志记录**：详细记录操作日志和错误日志，便于问题排查和系统维护。
   - **异常处理**：处理获取数据和计算排名过程中可能出现的各种异常情况，确保系统的稳定性和可靠性。
3. **性能优化**：
   - **批量处理**：通过批量处理帖子数据和交互数据，减少数据库的访问次数，提高处理效率。
   - **优先队列优化**：在队列已满时，进行替换操作，确保队列中始终保留分数最高的帖子，提高计算效率。

## 评论模块

在评论模块的设计中，我采取了以下关键步骤来优化性能和用户体验：

1. **临接表设计**：
   - **多级评论支持**：通过设计根评论ID (`RootId`) 和父评论ID (`PID`) 两个字段来支持多级评论结构，确保评论层级关系的清晰。
   - **父评论引用**：通过 `ParentComment` 字段实现对父评论的引用，使得每条评论都可以关联其父评论，形成树状结构。
   - **高效查询**：使用索引优化根评论ID和父评论ID的查询，提高获取评论和回复的速度。

2. **数据访问层设计**：
   - **评论创建与删除**：通过 `CreateComment` 和 `DeleteComment` 方法实现评论的创建与删除，确保数据操作的高效性和一致性。
   - **评论查询**：
      - `ListComments` 方法实现分页查询指定帖子的评论，确保在大数据量时仍能高效获取评论列表。
      - `GetMoreCommentsReply` 方法实现分页获取更多评论回复，支持加载更多评论回复的功能，提升用户体验。

3. **并发处理**：
   - 使用 `errgroup` 并发获取每个评论的子评论，确保在不阻塞主线程的情况下，快速获取所有评论及其回复，提升查询效率。

4. **数据转换**：
   - **DAO 与领域模型的转换**：
      - `toDAOComment` 方法实现领域模型评论到 DAO 评论的转换，确保数据在不同层次间的正确传递。
      - `toDomainComment` 方法实现 DAO 评论到领域模型评论的转换，确保业务逻辑处理的数据类型一致性。
      - `toDomainSliceComments` 方法批量转换 DAO 评论为领域模型评论，提升数据转换的效率。

5. **安全性与性能**：
   - **日志记录**：使用 `zap` 日志库记录关键操作和错误，便于调试和监控。
   - **异常处理**：在数据操作和并发任务中处理可能出现的错误，确保系统的稳定性和可靠性。


### 技术优势

1. **临接表设计**：
   - **多级评论支持**：通过根评论ID和父评论ID的设计，支持多层次的评论结构，确保评论关系的清晰和查询的高效。
   - **索引优化**：对根评论ID和父评论ID添加索引，提高查询性能，特别是在大量评论数据的情况下。

2. **并发处理**：
   - **errgroup 并发处理**：通过 `errgroup` 并发获取子评论，提升查询效率，确保系统的高性能。
   - **批量处理**：通过分页加载和批量转换数据，减少数据库的访问次数，提高处理效率。

3. **性能优化**：
   - **批量处理**：通过分页加载和批量转换数据，减少数据库的访问次数，提高处理效率。
   - **数据转换优化**：通过高效的数据转换方法，确保数据在不同层次间的快速转换和传递。

## 搜索模块

在搜索模块的设计中，采取了以下关键步骤来优化性能和用户体验：

1. **Elasticsearch 集成**：
   - **高效搜索**：使用 Elasticsearch 提供全文搜索和关键字搜索功能，确保在海量数据下依然能快速响应用户的搜索请求。
   - **多索引支持**：通过创建 `PostIndex` 和 `UserIndex` 两个索引，分别存储帖子和用户的数据，便于分类管理和高效搜索。

2. **搜索功能实现**：
   - **搜索帖子**：`SearchPosts` 方法通过关键字搜索帖子内容和标题，并使用布尔查询条件确保搜索结果包含已发布状态的帖子，提升搜索结果的相关性和准确性。
   - **搜索用户**：`SearchUsers` 方法通过关键字搜索用户的邮箱、昵称和电话，提供多维度的用户搜索功能，增强用户查找体验。

3. **数据管理**：
   - **数据输入**：通过 `InputUser` 和 `InputPost` 方法将用户和帖子数据索引到 Elasticsearch，确保数据的实时性和准确性。
   - **数据删除**：通过 `DeleteUserIndex` 和 `DeletePostIndex` 方法删除 Elasticsearch 索引中的用户和帖子数据，确保数据的一致性和有效管理。

4. **查询构建与解析**：
   - **查询构建**：构建符合 Elasticsearch 查询 DSL 的 JSON 请求，使用多匹配查询和布尔查询来满足复杂搜索需求。
   - **查询解析**：解析 Elasticsearch 返回的搜索结果，将 JSON 数据反序列化为 `PostSearch` 和 `UserSearch` 对象，确保数据格式的一致性和正确性。

5. **日志记录与错误处理**：
   - **详细日志记录**：使用 `zap` 日志库记录关键操作和错误，便于调试和监控。
   - **错误处理**：处理在搜索请求、数据输入和删除操作中可能出现的错误，确保系统的稳定性和可靠性。

### 技术优势

1. **Elasticsearch 优势**：
   - **全文搜索**：通过 Elasticsearch 的全文搜索功能，能够快速、准确地从海量数据中找到相关内容。
   - **分布式架构**：Elasticsearch 的分布式架构支持大规模数据的高效存储和搜索，保证系统的高可用性和扩展性。
   - **丰富的查询语法**：Elasticsearch 提供多种查询语法，支持复杂的搜索需求，提升用户搜索体验。

2. **高效数据管理**：
   - **索引管理**：通过 `InputUser` 和 `InputPost` 方法，将数据实时索引到 Elasticsearch，确保搜索数据的实时性。
   - **数据清理**：通过 `DeleteUserIndex` 和 `DeletePostIndex` 方法，及时清理无效数据，确保索引数据的一致性和有效性。

3. **安全性与性能**：
   - **日志记录**：详细记录操作日志和错误日志，便于问题排查和系统维护。
   - **错误处理**：处理在搜索请求、数据输入和删除操作中可能出现的各种异常情况，确保系统的稳定性和可靠性。

## 缓存击穿、穿透、雪崩优化

### 采取以下策略来优化缓存击穿、穿透、雪崩问题：
   - 布谷鸟缓存：使用布谷鸟缓存来避免缓存击穿问题，确保缓存数据的一致性和有效性。
   - 热点数据永不过期：对于热点数据，设置永不过期的缓存，避免缓存击穿问题。

## 缓存、mysql与mongo、es等数据一致性问题

### 采取以下策略来优化缓存数据一致性问题：
   - kafka+canal 实现数据同步：使用 kafka+canal 实现数据同步，确保缓存数据的一致性和有效性。
