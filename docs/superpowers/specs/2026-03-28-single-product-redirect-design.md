# 单结果自动跳转产品详情页

## 背景

首页（SolutionAtlas）中的方案卡片下方会显示厂商标签（如 SmartX）。点击后会跳转到 `/products?solution_category_id=19&manufacturer=SmartX` 搜索结果页。

即使该厂商在该解决方案下只有一个产品，用户也需要多点击一次才能看到详情页。

## 目标

当搜索结果**仅有1个产品**时，自动跳转到该产品的详情页，减少不必要的中间页面。

## 设计方案

**文件：** `frontend/src/app/(marketing)/products/page.tsx`

**改动点：** 在 `ProductsInner` 组件的 `load()` 函数中，`setList` 后增加单结果判断逻辑：

```tsx
const isFirstLoad = useRef(true);

const load = useCallback(async () => {
  setLoading(true);
  // ... existing load logic ...

  const res = await apiFetch(path);
  if (res.ok) {
    const data = await res.json();
    // 如果只有1个结果，自动跳转到详情页
    if (data.length === 1 && !isFirstLoad.current) {
      isFirstLoad.current = true;
      router.replace(`/products/${data[0].id}`);
      setLoading(false);
      return;
    }
    setList(data);
  } else {
    setList([]);
  }
  isFirstLoad.current = false;
  setLoading(false);
}, [/* deps */]);
```

## 行为说明

- **任意筛选条件组合**（厂商、解决方案、市场、关键词）导致单结果时均跳转
- 使用 `router.replace` 而非 `router.push`，避免产生浏览器历史记录堆叠
- `isFirstLoad` ref 防止首次加载时误触发跳转（首次加载由 URL 参数驱动，应正常显示搜索结果）
- 不影响多结果时的正常搜索列表展示

## 影响范围

- 仅修改 `products/page.tsx` 一个文件
- 不改变任何后端接口逻辑
- 不影响其他页面跳转行为
