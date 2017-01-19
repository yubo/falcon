  </div>
 </div>
</div>
<div id="rMenu" class="btn-group">
  <ul class="dropdown-menu">
    <li id="m_add" onclick="addTreeNode();">增加节点</li>
    <li id="m_del" onclick="removeTreeNode();">删除节点</li>
    <li id="m_check" onclick="checkTreeNode(true);">Check节点</li>
    <li id="m_unCheck" onclick="checkTreeNode(false);">unCheck节点</li>
    <li id="m_reset" onclick="resetTree();">恢复zTree</li>
  </ul>
</div>
{{template "layout/foot.tpl" .}}
