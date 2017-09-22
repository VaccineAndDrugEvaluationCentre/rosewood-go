Private Sub Document_Open()
'Preferred zoom %
Const lZoom As Long = 120
With ActiveWindow
  'Ignore any errors
  On Error Resume Next
  'Reduce flickering while changing settings
  .Visible = False
  'Switch to a single view pane
  .View.SplitSpecial = wdPaneNone
  'Switch to Normal mode
  .View.Type = wdNormalView
  With .ActivePane.View.Zoom
    'Initialize for best fit
    .PageFit = wdPageFitBestFit
    'Round down to nearest 10%
    .Percentage = Int(.Percentage / 10) * 10
    'Test zoom % and reduce to lZoom% max
    If .Percentage > lZoom Then .Percentage = lZoom
  End With
  'Switch to Print Preview mode
  .View.Type = wdPrintView
  With .ActivePane.View.Zoom
    'Set for a 1-page view
    .PageColumns = 1
    'Initialize for best fit
    .PageFit = wdPageFitBestFit
    'Round down to nearest 10%
    .Percentage = Int(.Percentage / 10) * 10
    'Test zoom % and reduce to lZoom% max
    If .Percentage > lZoom Then .Percentage = lZoom
  End With
  'Display the Rulers
  .ActivePane.DisplayRulers = True
  'Restore the window now that we're finished
  .Visible = True
End With
End Sub
